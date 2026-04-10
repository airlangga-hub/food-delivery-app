# Food Delivery
Group two application for Hacktiv8 phase three final project.

## 👥 Team Members

| Name                                                               | 
| ------------------------------------------------------------------ |
| [Airlangga Krishna](https://github.com/airlangga-hub)              |
| [Edric Emerson](https://github.com/edricemerson)                   | 
| [Mohammad Firmansyah Suryo Baskoro](https://github.com/Firmeteran) |

## Technologies Used

| Technology                                     | Description                                                                                          |
| ---------------------------------------------- | ---------------------------------------------------------------------------------------------------- |
| [Go (programming language)](https://go.dev/)   | Main programming language, known for efficiency and concurrency.                                     |
| [Echo](https://echo.labstack.com/)             | Minimalist web framework for building high-performance REST APIs.                                    |
| [gRPC](https://grpc.io/)                       | High performance, open source universal RPC framework.                                               |
| [Supabase](https://supabase.com/)              | Cloud-based PostgreSQL platform for reliable data storage.                                           |
| [Xendit](https://www.xendit.co)                | Payment gateway integration for automated transaction handling.                                      |
| [Mailjet](https://www.mailjet.com/)            | Cloud-based solution for designing, sending, and tracking both marketing and transactional emails.   |

## Project Structure

```
food-delivery-app/
├── gateway/
│   ├── auth/
│   │   └── auth.go            
│   ├── handler/
│   │   ├── entity.go
│   │   └── handler.go
│   ├── helper/
│   │   └── jwt.go  
│   ├── model/
│   │   ├── const.go
│   │   ├── error.go
│   │   ├── order.go
│   │   └── user.go
│   ├── order_pb/
│   │   ├── order_grpc.pb.go
│   │   ├── order.pb.go
│   │   └── order.proto
│   ├── service/
│   │   ├── order_service.go
│   │   └── user_service.go
│   ├── user_pb/
│   │   ├── user_grpc.pb.go
│   │   ├── user.pb.go
│   │   └── user.proto
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   └── main.go
├── order/                  
│   ├── auth/
│   │   └── auth.go            
│   ├── handler/
│   │   └── grpc.go
│   ├── helper/
│   │   └── ptr.go  
│   ├── middleware/
│   │   └── middleware.go
│   ├── model/
│   │   ├── const.go
│   │   ├── error.go
│   │   ├── order.go
│   │   ├── payment_gateway.go
│   │   └── payment_record.go
│   ├── pb/
│   │   ├── order_grpc.pb.go
│   │   ├── order.pb.go
│   │   └── order.proto
│   ├── repository/
│   │   ├── mongo.go
│   │   ├── sql_entity.go
│   │   ├── sql.go
│   │   ├── xendit_entity.go
│   │   └── xendit.go
│   ├── service/
│   │   ├── customer.go
│   │   └── driver.go
│   ├── user_pb/
│   │   ├── user_grpc.pb.go
│   │   ├── user.pb.go
│   │   └── user.proto
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   └── main.go           
├── user/                 
│   ├── auth/
│   │   └── auth.go            
│   ├── handler/
│   │   └── grpc.go
│   ├── helper/
│   │   └── jwt.go
│   ├── middleware/
│   │   └── middleware.go  
│   ├── model/
│   │   ├── const.go
│   │   ├── error.go
│   │   ├── mailjet.go
│   │   ├── payment_gateway.go
│   │   ├── payment_record.go
│   │   └── user.go
│   ├── order_pb/
│   │   ├── order_grpc.pb.go
│   │   ├── order.pb.go
│   │   └── order.proto
│   ├── pb/
│   │   ├── user_grpc.pb.go
│   │   ├── user.pb.go
│   │   └── user.proto
│   ├── repository/
│   │   ├── mongo.go
│   │   ├── payment_gateway.go
│   │   └── sql.go
│   ├── service/
│   │   └── user.go
│   ├── user_pb/
│   │   └── database/
│   │       └── database.go
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   └── main.go        
├── .env.example   
├── .gitignore
├── docker-compose.local.yaml      
├── docker-compose.yaml      
└── Makefile                
```

## Installation & How to Run
1. Clone the repository
```
git clone https://github.com/airlangga-hub/food-delivery-app.git
cd food-delivery-app
```

2. Install dependencies
```
go mod tidy
```

3. Run the application
```
go run gateway/main.go
```
