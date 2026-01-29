# gRPC CRUD Project

A gRPC-based CRUD application with REST API gateway for user management.

## Features

- gRPC service for User CRUD operations
- HTTP/REST gateway for Postman/browser access
- PostgreSQL database integration
- Protocol Buffers for type-safe communication

## Prerequisites

- Go 1.21+
- PostgreSQL
- Protocol Buffers compiler (`protoc`)

## Setup

1. Install dependencies:
```bash
chmod +x setup.sh
./setup.sh
```

2. Create database table:
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE
);
```

3. Run the server:
```bash
go run server/main.go
```

## API Endpoints

- `POST /v1/users` - Create user
- `GET /v1/users/{id}` - Get user
- `PUT /v1/users/{id}` - Update user
- `DELETE /v1/users/{id}` - Delete user

## Project Structure

```
grpc-crud-proj/
├── proto/          # Protocol buffer definitions
├── server/         # gRPC server implementation
├── client/         # gRPC client example
└── db/             # Database connection
```

## Testing

Use the provided Postman collection or test with curl:

```bash
# Create user
curl -X POST http://localhost:8080/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"John","email":"john@example.com"}'

# Get user
curl http://localhost:8080/v1/users/1

# Update user
curl -X PUT http://localhost:8080/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"John Updated","email":"john.updated@example.com"}'

# Delete user
curl -X DELETE http://localhost:8080/v1/users/1
```
