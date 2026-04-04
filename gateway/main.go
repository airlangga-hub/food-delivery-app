package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/airlangga-hub/food-delivery-app/gateway/auth"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	godotenv.Load()

	port := os.Getenv("CONTAINER_PORT")
	jwtSecret := os.Getenv("JWT_SECRET")
	userAddress := os.Getenv("USER_ADDRESS")
	orderAddress := os.Getenv("ORDER_ADDRESS")
	grpcUsername := os.Getenv("GRPC_USERNAME")
	grpcPassword := os.Getenv("GRPC_PASSWORD")

	if port == "" || jwtSecret == "" || orderAddress == "" || userAddress == "" || grpcUsername == "" || grpcPassword == "" {
		log.Fatalln("env variable missing.")
	}

	// grpc basic auth
	auth := auth.BasicAuth{
		Username: grpcUsername,
		Password: grpcPassword,
	}

	// user grpc client conn
	userCC, err := grpc.NewClient(userAddress, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithPerRPCCredentials(auth))
	if err != nil {
		logger.Error("user gRPC client conn error", slog.Any("error", err))
		return
	}
	defer userCC.Close()

	// order grpc client conn
	orderCC, err := grpc.NewClient(orderAddress, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithPerRPCCredentials(auth))
	if err != nil {
		logger.Error("order gRPC client conn error", slog.Any("error", err))
		return
	}
	defer orderCC.Close()

	// dependency injection
	userClient := userpb.NewUserServiceClient(userCC)
	orderClient := orderpb.NewOrderServiceClient(orderCC)
	svc := service.New(userClient, orderClient)
	validate := validator.New(validator.WithRequiredStructEnabled())
	h := handler.New(svc, validate)

	// echo
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:    true,
		LogMethod:    true,
		LogURI:       true,
		LogRequestID: true,
		LogLatency:   true,
		HandleError:  true, // forwards error to the global error handler, so it can decide appropriate status code
		LogValuesFunc: func(c *echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				logger.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST",
					slog.Int("status", v.Status),
					slog.String("method", v.Method),
					slog.String("uri", v.URI),
					slog.String("request_id", v.RequestID),
					slog.Int64("latency", int64(v.Latency)),
				)
			} else {
				logger.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR",
					slog.Int("status", v.Status),
					slog.String("method", v.Method),
					slog.String("uri", v.URI),
					slog.String("request_id", v.RequestID),
					slog.Int64("latency", int64(v.Latency)),
					slog.String("error", v.Error.Error()),
				)
			}
			return nil
		},
	}))
}
