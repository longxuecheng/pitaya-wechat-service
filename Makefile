BINARY_NAME=gotrue
DEST_DIR=/usr/pitaya/
APP_SERVER=root@vm1

build: VER=$(shell git rev-parse --short HEAD)
build:
	go build -o $(BINARY_NAME) -v

clean:
	go clean

build-linux: clean
	CGO_ENABLED=0 GOOS=linux make build

deploy: build-linux
	scp $(BINARY_NAME) $(APP_SERVER):$(DEST_DIR)