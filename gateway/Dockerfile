# Stage 1: Build the Go application binary
FROM golang:1.25.3 AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o gateway

# Stage 2: Create a minimal runtime image and running the application
FROM debian:bookworm-slim
RUN apt-get update && \
      apt-get install -y --no-install-recommends ca-certificates && \
      rm -rf /var/lib/apt/lists/*
COPY --from=build /app/gateway /gateway

ENTRYPOINT ["/gateway"]