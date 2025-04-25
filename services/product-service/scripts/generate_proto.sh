#!/bin/bash

# Ensure the script directory exists
mkdir -p scripts

# Generate gRPC code from proto files
protoc \
  --go_out=. \
  --go_opt=paths=source_relative \
  --go-grpc_out=. \
  --go-grpc_opt=paths=source_relative \
  proto/product/product.proto

echo "Proto generation completed!" 