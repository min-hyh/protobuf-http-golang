#!/bin/bash

echo "Starting Discover Service HTTP Server..."
echo "========================================"

# Check if the protobuf files exist
if [ ! -f "pb/discover.pb.go" ]; then
    echo "Error: Protobuf files not found. Please run the protoc command first:"
    echo "protoc -I. -I./third_party -I./pb --go_out=. --go-grpc_out=. --grpc-gateway_out=. pb/discover.proto"
    exit 1
fi

# Run the server
echo "Starting server..."
go run server/main.go 