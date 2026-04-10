gen-user:
	cd user/pb && \
		protoc \
			--go_out=. --go_opt=paths=source_relative \
			--go-grpc_out=. --go-grpc_opt=paths=source_relative \
			user.proto

gen-order:
	cd order/pb && \
		protoc \
			--go_out=. --go_opt=paths=source_relative \
			--go-grpc_out=. --go-grpc_opt=paths=source_relative \
			order.proto

logs:
	docker compose logs -f

local-build:
	docker compose -f docker-compose.yaml -f docker-compose.local.yaml up --build -d

local:
	docker compose -f docker-compose.yaml -f docker-compose.local.yaml up -d

build-gateway:
	docker build --platform linux/amd64 -t airlangga491/final-project-gateway:latest ./gateway/

push-gateway:
	docker push airlangga491/final-project-gateway:latest

build-order:
	docker build --platform linux/amd64 -t airlangga491/final-project-order:latest ./order/

push-order:
	docker push airlangga491/final-project-order:latest

build-user:
	docker build --platform linux/amd64 -t airlangga491/final-project-user:latest ./user/

push-user:
	docker push airlangga491/final-project-user:latest