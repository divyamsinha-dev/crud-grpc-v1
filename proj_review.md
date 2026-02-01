# gRPC CRUD Project Review

## Project Overview

A gRPC-based CRUD application with REST API gateway for user management.

### Features

- gRPC service for User CRUD operations
- HTTP/REST gateway for Postman/browser access
- PostgreSQL database integration
- Protocol Buffers for type-safe communication

### API Endpoints

- `POST /v1/users` - Create user
- `GET /v1/users/{id}` - Get user
- `PUT /v1/users/{id}` - Update user
- `DELETE /v1/users/{id}` - Delete user

### Project Structure

```
grpc-crud-proj/
├── proto/          # Protocol buffer definitions
├── server/         # gRPC server implementation
├── client/         # gRPC client example
└── db/             # Database connection
```

---

## Deployment Guide (Docker & Kubernetes)

This section explains how to containerize the application and deploy it to a Kubernetes cluster.

### 1. Dockerization

The project includes a `Dockerfile` for multi-stage builds (using `golang:alpine` for building and `alpine` for running).

**Build the Image:**
```bash
docker build -t grpc-crud-app:latest .
```

### 2. Kubernetes Deployment

The `k8s/` folder contains manifests for the application and PostgreSQL.

**Manifests:**
- **`postgres-configmap.yaml`**: Database initialization SQL.
- **`postgres-deployment.yaml`**: runs PostgreSQL.
- **`postgres-service.yaml`**: Exposes DB at `postgres:5432`.
- **`app-deployment.yaml`**: Runs the Go app.
- **`app-service.yaml`**: Exposes app ports (8080/50051).

### 3. Deploying

**Load Image (Minikube/Kind only):**
```bash
minikube image load grpc-crud-app:latest
# OR
kind load docker-image grpc-crud-app:latest
```

**Apply Manifests:**
```bash
# Deploy Database
kubectl apply -f k8s/postgres-configmap.yaml
kubectl apply -f k8s/postgres-deployment.yaml
kubectl apply -f k8s/postgres-service.yaml

# Deploy Application
kubectl apply -f k8s/app-deployment.yaml
kubectl apply -f k8s/app-service.yaml
```

**Access Application:**
If using Minikube, run `minikube tunnel`. Then access at `localhost:8080`.

### Testing

```bash
# Create user
curl -X POST http://localhost:8080/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"John","email":"john@example.com"}'
```
