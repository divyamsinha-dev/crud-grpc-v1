package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	"grpc-crud-proj/db"
	gw "grpc-crud-proj/proto/google/userpb"
	pb "grpc-crud-proj/proto/google/userpb"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type server struct {
	pb.UnimplementedUserServiceServer
	db *sql.DB
}

// Add this inside server/main.go

func (s *server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.UserResponse, error) {
	hashedPwd, _ := hashPassword(req.Password)

	// Default to "user" if no role is sent
	userRole := req.Role
	if userRole == "" {
		userRole = "user"
	}

	var id int
	// INSERT the role into DB
	err := s.db.QueryRow(
		"INSERT INTO users(name, email, password, role) VALUES($1, $2, $3, $4) RETURNING id",
		req.Name, req.Email, hashedPwd, userRole,
	).Scan(&id)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create user: %v", err)
	}

	return &pb.UserResponse{
		User: &pb.User{Id: int32(id), Name: req.Name, Email: req.Email, Role: userRole},
	}, nil
}

func (s *server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	var storedHash string
	var role string // <--- 1. Variable to hold the role

	// 2. CRITICAL: We must SELECT the 'role' column from the DB
	err := s.db.QueryRow(
		"SELECT password, role FROM users WHERE email=$1",
		req.Email,
	).Scan(&storedHash, &role) // <--- 3. Scan it into the variable

	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "user not found")
	}

	if !checkPassword(req.Password, storedHash) {
		return nil, status.Errorf(codes.Unauthenticated, "incorrect password")
	}

	// 4. Pass the fetched role to the token generator
	token, err := generateToken(req.Email, role)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot generate token")
	}

	return &pb.LoginResponse{Token: token}, nil
}

func (s *server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	var id int
	// Include the role in the INSERT statement
	err := s.db.QueryRow(
		"INSERT INTO users(name, email, role) VALUES($1, $2, $3) RETURNING id",
		req.Name, req.Email, req.Role,
	).Scan(&id)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	return &pb.UserResponse{
		User: &pb.User{
			Id:    int32(id),
			Name:  req.Name,
			Email: req.Email,
			Role:  req.Role,
		},
	}, nil
}

func (s *server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	var user pb.User
	// Add 'role' to the SELECT and Scan
	err := s.db.QueryRow(
		"SELECT id, name, email, role FROM users WHERE id=$1",
		req.Id,
	).Scan(&user.Id, &user.Name, &user.Email, &user.Role)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, err
	}

	return &pb.UserResponse{User: &user}, nil
}

func (s *server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	_, err := s.db.Exec(
		"UPDATE users SET name=$1, email=$2 WHERE id=$3",
		req.Name, req.Email, req.Id,
	)
	if err != nil {
		return nil, err
	}

	return &pb.UserResponse{
		User: &pb.User{
			Id:    req.Id,
			Name:  req.Name,
			Email: req.Email,
		},
	}, nil
}

func (s *server) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	_, err := s.db.Exec("DELETE FROM users WHERE id=$1", req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteUserResponse{
		Message: "User deleted",
	}, nil
}

func main() {
	dbConn := db.Connect()

	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatal("Failed to listen on gRPC port:", err)
		}

		//grpcServer := grpc.NewServer()
		// We register the interceptor here!
		grpcServer := grpc.NewServer(
			grpc.UnaryInterceptor(AuthInterceptor),
		)
		pb.RegisterUserServiceServer(grpcServer, &server{db: dbConn})

		log.Println("gRPC server running on :50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal("Failed to serve gRPC:", err)
		}
	}()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal("Failed to dial gRPC server:", err)
	}
	defer conn.Close()

	mux := runtime.NewServeMux()

	err = gw.RegisterUserServiceHandler(ctx, mux, conn)
	if err != nil {
		log.Fatal("Failed to register gateway:", err)
	}

	log.Println("HTTP/REST gateway running on :8080")
	log.Println("POST   http://localhost:8080/v1/users")
	log.Println("GET    http://localhost:8080/v1/users/{id}")
	log.Println("PUT    http://localhost:8080/v1/users/{id}")
	log.Println("DELETE http://localhost:8080/v1/users/{id}")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal("Failed to serve HTTP:", err)
	}
}
