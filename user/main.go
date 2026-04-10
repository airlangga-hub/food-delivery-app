package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"

	"github.com/airlangga-hub/food-delivery-app/user/auth"
	"github.com/airlangga-hub/food-delivery-app/user/handler"
	"github.com/airlangga-hub/food-delivery-app/user/middleware"
	orderpb "github.com/airlangga-hub/food-delivery-app/user/order_pb"
	"github.com/airlangga-hub/food-delivery-app/user/pb"
	"github.com/airlangga-hub/food-delivery-app/user/repository"
	"github.com/airlangga-hub/food-delivery-app/user/service"
	"github.com/airlangga-hub/food-delivery-app/user/util/database"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	godotenv.Load()

	mongoURI := os.Getenv("MONGO_URI")
	port := os.Getenv("CONTAINER_PORT")
	dbMongoName := os.Getenv("DB_MONGO_NAME")
	mailjetSender := os.Getenv("MAILJET_SENDER")
	mailjetURL := os.Getenv("MAILJET_URL")
	mailjetUsername := os.Getenv("MAILJET_BASIC_AUTH_USERNAME")
	mailjetPassword := os.Getenv("MAILJET_BASIC_AUTH_PASSWORD")
	grpcUsername := os.Getenv("GRPC_USERNAME")
	grpcPassword := os.Getenv("GRPC_PASSWORD")
	orderAddress := os.Getenv("ORDER_ADDRESS")
	supabaseURI := os.Getenv("SUPABASE_URI")
	jwtKey := os.Getenv("JWT_SECRET")

	if jwtKey == "" || supabaseURI == "" || grpcUsername == "" || grpcPassword == "" || mongoURI == "" || port == "" || dbMongoName == "" || mailjetSender == "" || mailjetURL == "" || mailjetUsername == "" || mailjetPassword == "" || orderAddress == "" {
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

	// grpc basic auth
	auth := auth.BasicAuth{
		Username: grpcUsername,
		Password: grpcPassword,
	}

	// order grpc client conn
	orderCC, err := grpc.NewClient(orderAddress, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithPerRPCCredentials(auth))
	if err != nil {
		logger.Error("order gRPC client conn error", slog.Any("error", err))
		return
	}
	defer orderCC.Close()

	// sql db
	sqlDB, err := sql.Open("postgres", supabaseURI)
	if err != nil {
		logger.Error("error connecting to supabase", slog.Any("error", err))
		return
	}
	defer sqlDB.Close()

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	if err := sqlDB.Ping(); err != nil {
		logger.Error("error pinging supabase", slog.Any("error", err))
		return
	}
	
	_, err = sqlDB.Exec("SET search_path TO final_project")
	if err != nil {
		logger.Error("error setting search path in supabase", slog.Any("error", err))
		return
	}

	// dependency injection
	orderClient := orderpb.NewOrderServiceClient(orderCC)
	validate := validator.New(validator.WithRequiredStructEnabled())
	userPaymentGatewayRepo := repository.NewPaymentGatewayRepository(orderClient)
	userMongoRepo := repository.NewMongoRepository(client.Database(dbMongoName), validate, mailjetSender, mailjetURL, mailjetUsername, mailjetPassword)
	userSQLRepo := repository.NewSQLRepository(sqlDB)
	svc := service.NewUserService(userPaymentGatewayRepo, userMongoRepo, userSQLRepo, logger, jwtKey)

	// cron
	c := cron.New(cron.WithSeconds())

	_, err = c.AddFunc("0 */1 * * * *", func() {
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
	basicAuthMap[grpcUsername] = grpcPassword

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
	pb.RegisterUserServiceServer(grpcServer, handler.NewHandler(svc))

	logger.Info("gRPC Server running on port " + port)

	chanExit := make(chan string)

	go func() {
		chanSig := make(chan os.Signal, 1)
		signal.Notify(chanSig, syscall.SIGINT, syscall.SIGTERM)
		s := <-chanSig
		chanExit <- fmt.Sprintf("Signal %v received, shutting down gRPC server...", s)
	}()

	go func() {
		lis, err := net.Listen("tcp", ":"+port)
		if err != nil {
			chanExit <- fmt.Sprintf("net.Listen error: %v", err)
			return
		}
		if err = grpcServer.Serve(lis); err != nil && err != grpc.ErrServerStopped {
			chanExit <- fmt.Sprintf("gRPC.Serve error: %v", err)
		}
	}()
	
	exitSignal := <-chanExit

	logger.Info(fmt.Sprintf("Exit signal received: %s", exitSignal))

	grpcServer.GracefulStop()

	logger.Info("Graceful shutdown complete.")
}
