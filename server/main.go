package main

import (
	"context"
	"database/sql"
	"log"
	"net"

	"grpc-crud-proj/db"
	pb "grpc-crud-proj/proto/userpb"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedUserServiceServer
	db *sql.DB
}

func (s *server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	var id int
	err := s.db.QueryRow(
		"INSERT INTO users(name,email) VALUES($1,$2) RETURNING id",
		req.Name, req.Email).Scan(&id)

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

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, &server{db: dbConn})

	log.Println("gRPC server running on :50051")
	grpcServer.Serve(lis)
}
