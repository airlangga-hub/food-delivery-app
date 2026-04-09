package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"github.com/airlangga-hub/food-delivery-app/user/handler"
	"github.com/airlangga-hub/food-delivery-app/user/middleware"
	"github.com/airlangga-hub/food-delivery-app/user/pb"
	"github.com/airlangga-hub/food-delivery-app/user/repository"
	"github.com/airlangga-hub/food-delivery-app/user/service"
	"github.com/airlangga-hub/food-delivery-app/user/util/database"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"google.golang.org/grpc"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	godotenv.Load()

	mongoURI := os.Getenv("MONGO_URI")
	port := os.Getenv("PORT")
	xenditAPIkey := os.Getenv("XENDIT_API_KEY")
	xenditPaymentSessionURL := os.Getenv("XENDIT_PAYMENT_SESSION_URL")
	dbMongoName := os.Getenv("DB_MONGO_NAME_PAYMENT")
	mailjetSender := os.Getenv("MAILJET_SENDER")
	mailjetURL := os.Getenv("MAILJET_URL")
	mailjetUsername := os.Getenv("MAILJET_BASIC_AUTH_USERNAME")
	mailjetPassword := os.Getenv("MAILJET_BASIC_AUTH_PASSWORD")
	grpcUser := os.Getenv("GRPC_USER")
	grpcPassword := os.Getenv("GRPC_PASSWORD")
	if grpcUser == "" || grpcPassword == "" || mongoURI == "" || port == "" || dbMongoName == "" || xenditAPIkey == "" || xenditPaymentSessionURL == "" || mailjetSender == "" || mailjetURL == "" || mailjetUsername == "" || mailjetPassword == "" {
		logger.Error("env variable missing.")
		return
	}

	client, err := database.GetMongoClient(mongoURI)
	if err != nil {
		logger.Error("get mongo client error", slog.Any("error", err))
		return
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			logger.Error("mongo disconnect error:", slog.Any("error", err))
		}
	}()

	// dependency injection
	validate := validator.New(validator.WithRequiredStructEnabled())
	repo, err := repository.NewRepository(client.Database(dbMongoName), xenditPaymentSessionURL, xenditAPIkey, mailjetSender, mailjetURL, mailjetUsername, mailjetPassword, validate)
	if err != nil {
		logger.Error("new mongo repo error", slog.Any("error", err))
		return
	}
	svc := service.NewUserService(repo, logger)

	// cron
	c := cron.New(cron.WithSeconds())

	_, err = c.AddFunc("0 */1 * * * *", func() {
		logger.Info("Cron: Starting email cronjob...")
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		if err := svc.ProcessUnsentEmail(ctx); err != nil {
			logger.Error("Cron Error", slog.Any("error", err))
		}
	})

	if err != nil {
		logger.Error("Failed to add func to cron", slog.Any("error", err))
	}

	// start cron
	c.Start()
	defer c.Stop()

	// grpc basic auth
	basicAuthMap := make(map[string]string)
	basicAuthMap[grpcUser] = grpcPassword

	// grpc server + middleware
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.LoggingInterceptor(logger),
			middleware.BasicAuthUnaryInterceptor(basicAuthMap),
		),
		grpc.ChainStreamInterceptor(
			middleware.BasicAuthStreamInterceptor(basicAuthMap),
		),
	)

	// Register the service implementation
	pb.RegisterPaymentServiceServer(grpcServer, handler.NewPaymentServer(svc))

	logger.Info("gRPC Server running on port " + port)

	chanExit := make(chan error)

	go func() {
		chanSig := make(chan os.Signal, 1)
		signal.Notify(chanSig, syscall.SIGINT, syscall.SIGTERM)
		s := <-chanSig
		chanExit <- fmt.Errorf("Signal %v received, shutting down gRPC server...", s)
	}()

	go func() {
		lis, err := net.Listen("tcp", ":"+port)
		if err != nil {
			chanExit <- fmt.Errorf("net.Listen error: %v", err)
			return
		}
		if err = grpcServer.Serve(lis); err != nil && err != grpc.ErrServerStopped {
			chanExit <- fmt.Errorf("gRPC.Serve error: %v", err)
		}
	}()

	err = <-chanExit

	logger.Error(fmt.Sprintf("Exit signal received: %v", err))

	grpcServer.GracefulStop()

	logger.Info("Graceful shutdown complete.")
}