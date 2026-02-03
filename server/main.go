package main

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"grpc-crud-proj/db"
	gw "grpc-crud-proj/proto/userpb"
	pb "grpc-crud-proj/proto/userpb"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type server struct {
	pb.UnimplementedUserServiceServer
	db *sql.DB
}

func (s *server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	//ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	//defer cancel()

	var id int
	fmt.Println(ctx.Value("user"))
	err := s.db.QueryRowContext(ctx,
		"INSERT INTO users(name,email) VALUES($1,$2) RETURNING id",
		req.Name, req.Email,
	).Scan(&id)

	if err != nil {
		log.Printf(fmt.Sprint("Got this err" + err.Error()))
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
	//ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	//defer cancel()

	time.Sleep(10 * time.Second)

	var user pb.User

	err := s.db.QueryRowContext(ctx,
		"SELECT id, name, email FROM users WHERE id=$1",
		req.Id,
	).Scan(&user.Id, &user.Name, &user.Email)

	if err != nil {
		return nil, err
	}

	return &pb.UserResponse{User: &user}, nil
}

func (s *server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	//	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	//defer cancel()

	_, err := s.db.ExecContext(ctx,
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

/*
func (s *server) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	select {
	case <-time.After(5 * time.Minute):

	case <-ctx.Done():
		log.Println(ctx.Err())
		return nil, ctx.Err()
	}

	_, err := s.db.ExecContext(ctx, "DELETE FROM users WHERE id=$1", req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteUserResponse{
		Message: "User deleted",
	}, nil
}


*/

func (s *server) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {

	//ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	//defer cancel()

	_, err := s.db.ExecContext(ctx, "DELETE FROM users WHERE id=$1", req.Id)
	if err != nil {

		return nil, err
	}

	return &pb.DeleteUserResponse{
		Message: "User deleted successfully",
	}, nil
}

func unaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	// Basic auth expects format: "Basic <base64(username:password)>"
	authStr := authHeader[0]
	if !strings.HasPrefix(authStr, "Basic ") {
		return nil, status.Errorf(codes.Unauthenticated, "invalid authorization format")
	}

	payload, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(authStr, "Basic "))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid base64 in auth")
	}

	pair := strings.SplitN(string(payload), ":", 2)
	if len(pair) != 2 || pair[0] != "admin" || pair[1] != "password" {
		return nil, status.Errorf(codes.Unauthenticated, "invalid username or password")
	}
	ctx = context.WithValue(ctx, "user", pair[0])

	// Timing
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	start := time.Now()

	type result struct {
		resp interface{}
		err  error
	}
	resChan := make(chan result, 1)

	// Run the handler in a goroutine
	go func() {
		resp, err := handler(ctx, req)
		resChan <- result{resp: resp, err: err}
	}()

	// Race
	select {
	case <-ctx.Done():
		log.Printf("Method %s TIMEOUT after 2s", info.FullMethod)
		return nil, status.Errorf(codes.DeadlineExceeded, "request took too long")

	case res := <-resChan:
		log.Printf("Method: %s, Duration: %s", info.FullMethod, time.Since(start))
		return res.resp, res.err
	}
}

func main() {
	dbConn := db.Connect()

	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatal("Failed to listen on gRPC port:", err)
		}

		grpcServer := grpc.NewServer(
			grpc.UnaryInterceptor(unaryInterceptor),
		)
		pb.RegisterUserServiceServer(grpcServer, &server{db: dbConn})

		log.Println("gRPC server running on :50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal("Failed to serve gRPC:", err)
		}
	}()

	ctx := context.Background()
	//ctx, cancel := context.WithCancel(ctx)
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
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

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal("Failed to serve HTTP:", err)
	}
}
