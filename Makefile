server-linux-build:
	GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/server ./cmd/server

client-linux-build:
	GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/client ./cmd/client

start-server: server-linux-build 
	docker-compose up --build server

start-client: client-linux-build
	docker-compose up --build client

test:
	go test ./...