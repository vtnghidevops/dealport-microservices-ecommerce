# Product Image Serving via gRPC

## Overview

This update modifies the product image serving mechanism in the ecommerce microservices project to use gRPC instead of direct HTTP. Previously, the broker service was directly proxying HTTP requests to the product service's HTTP endpoints to serve images. With this update, the broker service will now retrieve image data through gRPC and serve it to clients.

## Changes Made

1. Added new protobuf messages and service method:
   - `GetProductImageFileRequest`: Contains the image filename to retrieve
   - `GetProductImageFileResponse`: Contains the image binary data and content type
   - Added `GetProductImageFile` method to the ProductService service

2. Implemented the new gRPC service method in product-service:
   - Added implementation of `GetProductImageFile` that reads image files from the filesystem
   - Returns the image data and appropriate content type based on the file extension

3. Updated the broker-service to use gRPC for image retrieval:
   - Added `GetProductImageFile` method to the ProductClient in broker-service
   - Modified the `ProxyProductImage` handler to use gRPC instead of HTTP proxying
   - Added caching headers to improve performance

## Benefits

- **Consistent Communication Pattern**: All service-to-service communication now uses gRPC
- **Improved Security**: Direct filesystem access is limited to the product-service only
- **Better Performance**: gRPC is more efficient than HTTP for service-to-service communication
- **Added Caching**: New implementation includes cache control headers
- **Simpler Configuration**: No need to know the HTTP endpoint of the product service

## Required Steps to Implement

1. Generate the updated protobuf files for both services:
   ```bash
   # In product-service directory
   protoc --go_out=. --go-grpc_out=. proto/product/product.proto
   
   # In broker-service directory
   protoc --go_out=. --go-grpc_out=. proto/product/product.proto
   ```

2. Rebuild and restart both services:
   ```bash
   # Rebuild and restart product-service
   cd services/product-service
   go build -o product-service ./cmd/api
   ./product-service
   
   # Rebuild and restart broker-service
   cd services/broker-service
   go build -o broker-service ./cmd/api
   ./broker-service
   ```

## Note

This change eliminates the need for the direct HTTP route in the product-service for serving images to the broker. However, the existing HTTP handler in product-service is kept for backward compatibility or direct access when needed. 