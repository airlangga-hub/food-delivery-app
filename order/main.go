package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"

	"github.com/airlangga-hub/food-delivery-app/order/auth"
	"github.com/airlangga-hub/food-delivery-app/order/handler"
	"github.com/airlangga-hub/food-delivery-app/order/middleware"
	"github.com/airlangga-hub/food-delivery-app/order/order_pb"
	"github.com/airlangga-hub/food-delivery-app/order/repository"
	"github.com/airlangga-hub/food-delivery-app/order/service"
	userpb "github.com/airlangga-hub/food-delivery-app/order/user_pb"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	godotenv.Load()

	mongoURI := os.Getenv("MONGO_URI")
	port := os.Getenv("CONTAINER_PORT")
	xenditAPIkey := os.Getenv("XENDIT_API_KEY")
	xenditPaymentSessionURL := os.Getenv("XENDIT_PAYMENT_SESSION_URL")
	grpcUsername := os.Getenv("GRPC_USERNAME")
	grpcPassword := os.Getenv("GRPC_PASSWORD")
	userAddress := os.Getenv("USER_ADDRESS")
	supabaseURI := os.Getenv("SUPABASE_URI")

	if xenditAPIkey == "" || xenditPaymentSessionURL == "" || supabaseURI == "" || grpcUsername == "" || grpcPassword == "" || mongoURI == "" || port == "" || userAddress == "" {
		logger.Error("env variable missing.")
		return
	}

	// grpc basic auth
	auth := auth.BasicAuth{
		Username: grpcUsername,
		Password: grpcPassword,
	}

	// user grpc client conn
	userCC, err := grpc.NewClient(userAddress, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithPerRPCCredentials(auth))
	if err != nil {
		logger.Error("order gRPC client conn error", slog.Any("error", err))
		return
	}
	defer userCC.Close()

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

	// dependency injection
	userClient := userpb.NewUserServiceClient(userCC)
	validate := validator.New(validator.WithRequiredStructEnabled())
	sqlRepo := repository.NewSQLRepository(sqlDB)
	mongoRepo := repository.NewMongoRepository(userClient)
	xenditRepo := repository.NewXenditRepository(xenditPaymentSessionURL, xenditAPIkey, validate)
	customerSvc := service.NewCustomerService(sqlRepo, xenditRepo, mongoRepo)
	driverSvc := service.NewDriverService(sqlRepo)

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
	pb.RegisterOrderServiceServer(grpcServer, handler.NewHandler(customerSvc, driverSvc))

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
