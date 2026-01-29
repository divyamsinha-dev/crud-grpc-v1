package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	"grpc-crud-proj/db"
	gw "grpc-crud-proj/proto/userpb"
	pb "grpc-crud-proj/proto/userpb"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type server struct {
	pb.UnimplementedUserServiceServer
	db *sql.DB
}

func (s *server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	var id int
	err := s.db.QueryRow(
		"INSERT INTO users(name,email) VALUES($1,$2) RETURNING id",
		req.Name, req.Email,
	).Scan(&id)

	if err != nil {
		return nil, err
	}

	return &pb.UserResponse{
		User: &pb.User{
			Id:    int32(id),
			Name:  req.Name,
			Email: req.Email,
		},
	}, nil
}

func (s *server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	var user pb.User
	err := s.db.QueryRow(
		"SELECT id, name, email FROM users WHERE id=$1",
		req.Id,
	).Scan(&user.Id, &user.Name, &user.Email)

	if err != nil {
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

		grpcServer := grpc.NewServer()
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
