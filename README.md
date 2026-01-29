# gRPC CRUD Project with REST Gateway

A complete gRPC-based CRUD application that can be accessed via both gRPC clients and REST/HTTP (Postman).

## 🎯 What This Project Does

This project demonstrates:
- **gRPC** service implementation for User CRUD operations
- **gRPC Gateway** to expose REST endpoints for Postman/browser access
- **PostgreSQL** database integration
- **Complete learning guide** from beginner to advanced

## 🚀 Quick Start

### Prerequisites

1. **Go** (1.21+)
2. **PostgreSQL** (running locally)
3. **Protocol Buffers Compiler** (`protoc`)
   ```bash
   # macOS
   brew install protobuf
   
   # Linux
   apt-get install protobuf-compiler
   ```

### Setup Steps

1. **Install dependencies and generate code**:
   ```bash
   chmod +x setup.sh
   ./setup.sh
   ```

   This will:
   - Install Go dependencies
   - Generate gRPC and Gateway code from proto files

2. **Setup database**:
   ```sql
   CREATE TABLE users (
       id SERIAL PRIMARY KEY,
       name VARCHAR(255) NOT NULL,
       email VARCHAR(255) NOT NULL UNIQUE
   );
   ```

3. **Run the server**:
   ```bash
   go run server/main.go
   ```

   You should see:
   ```
   🚀 gRPC server running on :50051
   🌐 HTTP/REST gateway running on :8080
   ```

4. **Test with Postman**:

   **Create User**:
   - Method: `POST`
   - URL: `http://localhost:8080/v1/users`
   - Headers: `Content-Type: application/json`
   - Body:
     ```json
     {
       "name": "John Doe",
       "email": "john@example.com"
     }
     ```

   **Get User**:
   - Method: `GET`
   - URL: `http://localhost:8080/v1/users/1`

   **Update User**:
   - Method: `PUT`
   - URL: `http://localhost:8080/v1/users/1`
   - Headers: `Content-Type: application/json`
   - Body:
     ```json
     {
       "name": "John Updated",
       "email": "john.updated@example.com"
     }
     ```

   **Delete User**:
   - Method: `DELETE`
   - URL: `http://localhost:8080/v1/users/1`

5. **Test with gRPC client**:
   ```bash
   go run client/main.go
   ```

## 📁 Project Structure

```
grpc-crud-proj/
├── proto/
│   ├── user.proto              # Service definition
│   ├── google/api/              # HTTP annotations
│   └── userpb/                  # Generated code
│       ├── user.pb.go           # Message types
│       ├── user_grpc.pb.go      # gRPC service
│       └── user.pb.gw.go        # Gateway (generated)
├── server/
│   └── main.go                  # gRPC + Gateway server
├── client/
│   └── main.go                  # gRPC client example
├── db/
│   └── postgres.go              # Database connection
├── setup.sh                     # Setup script
├── GRPC_LEARNING_GUIDE.md       # Complete learning guide
└── README.md                     # This file
```

## 🎓 Learning Resources

**For Beginners - Start Here!**:
1. 📖 [`BEGINNER_WALKTHROUGH.md`](./BEGINNER_WALKTHROUGH.md) - Step-by-step explanation with simple analogies
2. 🎨 [`VISUAL_GUIDE.md`](./VISUAL_GUIDE.md) - Visual diagrams and flow charts
3. 💻 Read the code files - They have extensive comments explaining every line!

**For Advanced Learning**:
- 📚 [`GRPC_LEARNING_GUIDE.md`](./GRPC_LEARNING_GUIDE.md) - Complete guide from beginner to advanced

All code files now have detailed comments explaining:
- What each function does
- How the flow works
- Real-life analogies
- Step-by-step breakdowns

## 🔧 How It Works

### Architecture

```
┌─────────┐         ┌──────────────┐         ┌─────────────┐
│ Postman │ ─HTTP──>│ gRPC Gateway │ ─gRPC──>│ gRPC Server │
│ Browser │  JSON   │   (Port 8080)│ Protobuf│ (Port 50051)│
└─────────┘         └──────────────┘         └─────────────┘
                                                      │
                                                      ▼
                                               ┌─────────────┐
                                               │  PostgreSQL │
                                               └─────────────┘
```

### Flow

1. **Postman sends HTTP/JSON request** to `:8080`
2. **gRPC Gateway** translates JSON → Protobuf
3. **Gateway calls gRPC server** on `:50051`
4. **gRPC server** processes request (database operations)
5. **Response flows back** through the same path

## 🛠️ Troubleshooting

### "could not import grpc-gateway"
Run:
```bash
go mod tidy
go get github.com/grpc-ecosystem/grpc-gateway/v2/runtime
```

### "undefined: gw.RegisterUserServiceHandler"
Run the setup script to generate gateway code:
```bash
./setup.sh
```

### "protoc: command not found"
Install Protocol Buffers:
```bash
# macOS
brew install protobuf

# Linux
apt-get install protobuf-compiler
```

### Database connection issues
Check your PostgreSQL is running and update connection string in `db/postgres.go` or set `DB_URL` environment variable.

## 📝 Environment Variables

- `DB_URL`: PostgreSQL connection string (optional, defaults to localhost)

## 🎯 Key Features

✅ Full CRUD operations (Create, Read, Update, Delete)  
✅ Both gRPC and REST/HTTP access  
✅ Type-safe with Protocol Buffers  
✅ PostgreSQL integration  
✅ Complete learning guide  
✅ Production-ready structure  

## 📚 Next Steps

1. Read the [Learning Guide](./GRPC_LEARNING_GUIDE.md)
2. Try adding streaming RPCs
3. Add authentication/authorization
4. Implement error handling
5. Add logging and monitoring

## 🤝 Contributing

Feel free to extend this project with:
- Authentication
- More service types
- Streaming examples
- Error handling improvements
- Testing

## 📄 License

This is a learning project. Feel free to use and modify as needed.

---

**Happy Learning! 🚀**
