# Beginner's Guide to Dockerizing and Deploying to Kubernetes

This guide explains how to containerize your gRPC application and deploy it to a Kubernetes cluster.

## 1. Dockerization (The Container)

We created a `Dockerfile` in the root directory. This tells Docker how to build your application into an image.

### Concepts
- **Image**: A package containing your code, runtime, and system libraries.
- **Container**: A running instance of an image.
- **Multi-stage Build**: We use `golang:alpine` to *compile* the code (heavy image), but we run it in a plain `alpine` image (tiny image) to save space and improve security.

### Step 1: Build the Docker Image

Run this command in the project root:

```bash
docker build -t grpc-crud-app:latest .
```

This creates an image named `grpc-crud-app` with the tag `latest`.

---

## 2. Kubernetes (The Orchestrator)

We created a `k8s/` folder with several files. Kubernetes uses these "manifests" to manage your application.

### Concepts
- **Pod**: The smallest unit in K8s. Runs your container(s).
- **Deployment**: Manages Pods. Ensures a specific number of them are always running.
- **Service**: Exposes your Pods to the network. It gives them a stable IP and DNS name.
- **ConfigMap**: Stores configuration data (like our database initialization SQL).

### Files We Created
1. **`postgres-configmap.yaml`**: Contains the SQL to create the `users` table.
2. **`postgres-deployment.yaml`**: Runs the PostgreSQL database. It uses the ConfigMap to initialize the DB on startup.
3. **`postgres-service.yaml`**: Allows the app to find the database at `postgres:5432`.
4. **`app-deployment.yaml`**: Runs your Go application. It sets the `DB_URL` environment variable to connect to the database service.
5. **`app-service.yaml`**: Exposes your app to the outside world (Ports 8080 and 50051).

---

## 3. Deploying to Kubernetes

You can use Minikube, Kind, or Docker Desktop Kubernetes.

### Prerequisite: Load Image (If using Minikube/Kind)
Since we haven't pushed the image to a cloud registry (like Docker Hub), we need to load it into the cluster manually.

**For Minikube:**
```bash
minikube image load grpc-crud-app:latest
```

**For Kind:**
```bash
kind load docker-image grpc-crud-app:latest
```

**For Docker Desktop:**
It shares images with your local Docker, so you can skip this step!

### Step 2: Apply Manifests
Apply the database files first, then the application:

```bash
# 1. Deploy Database
kubectl apply -f k8s/postgres-configmap.yaml
kubectl apply -f k8s/postgres-deployment.yaml
kubectl apply -f k8s/postgres-service.yaml

# Wait a moment for DB to start...

# 2. Deploy Application
kubectl apply -f k8s/app-deployment.yaml
kubectl apply -f k8s/app-service.yaml
```

### Step 3: Check Status

Check if everything is running:
```bash
kubectl get pods
```
You should see `postgres-xxx` and `grpc-app-xxx` with status `Running`.

---

## 4. Accessing the Application

To access the service, you might need to tunnel if using Minikube:

```bash
minikube tunnel
```

Then you can access the HTTP Gateway at `localhost:8080`.

**Test with Curl:**
```bash
# Create a user
curl -X POST http://localhost:8080/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Docker User", "email":"docker@k8s.com"}'
```

If it works, you have successfully deployed a microservice architecture on Kubernetes!
