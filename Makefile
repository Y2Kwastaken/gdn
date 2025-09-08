.PHONY: build run clean

BINARY_NAME=grayva
CMD_DIR=./cmd/server

build:
	go build -o $(BINARY_NAME) $(CMD_DIR)
run:
	go run $(CMD_DIR)
test:
	go test -v ./...
clean:
	rm -f $(BINARY_NAME)