#!/bin/bash

# Setup script for gRPC Gateway
# This script installs dependencies and generates gateway code

echo "🔧 Setting up gRPC Gateway..."

# Install Go dependencies
echo "📦 Installing Go dependencies..."
go get github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
go get github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
go get google.golang.org/protobuf/cmd/protoc-gen-go
go get google.golang.org/grpc/cmd/protoc-gen-go-grpc
go get github.com/grpc-ecosystem/grpc-gateway/v2/runtime

# Check if protoc is installed
if ! command -v protoc &> /dev/null; then
    echo "❌ protoc is not installed. Please install Protocol Buffers compiler:"
    echo "   macOS: brew install protobuf"
    echo "   Linux: apt-get install protobuf-compiler"
    echo "   Or visit: https://grpc.io/docs/protoc-installation/"
    exit 1
fi

# Generate gRPC and Gateway code
echo "🔨 Generating gRPC and Gateway code..."
export PATH=$PATH:$(go env GOPATH)/bin
protoc \
  --proto_path=proto \
  --proto_path=. \
  --go_out=. \
  --go_opt=paths=source_relative \
  --go-grpc_out=. \
  --go-grpc_opt=paths=source_relative \
  --grpc-gateway_out=proto/userpb \
  --grpc-gateway_opt=paths=source_relative \
  proto/user.proto

if [ $? -eq 0 ]; then
    echo "✅ Code generation successful!"
else
    echo "❌ Code generation failed. Please check errors above."
    exit 1
fi

echo "✅ Setup complete!"
echo ""
echo "🚀 Next steps:"
echo "   1. Make sure PostgreSQL is running"
echo "   2. Create the users table: CREATE TABLE users (id SERIAL PRIMARY KEY, name VARCHAR(255), email VARCHAR(255));"
echo "   3. Run the server: go run server/main.go"
echo "   4. Test with Postman using the endpoints shown in the server logs"
