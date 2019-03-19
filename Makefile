BINARY_NAME=pitaya-wechat-service

build: VER=$(shell git rev-parse --short HEAD)
build:
	go build -o $(BINARY_NAME) -v

clean:
	go clean

build-linux: clean
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 make build
