package main

import (
	"context"
	"log"
	"time"

	pb "grpc-crud-proj/proto/userpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal("failed to connect:", err)
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	createRes, err := client.CreateUser(ctx, &pb.CreateUserRequest{
		Name:  "Divyam",
		Email: "divyam@test.com",
	})
	if err != nil {
		log.Fatal("CreateUser error:", err)
	}
	log.Println("Created User:", createRes.User)

	createRes2, err := client.CreateUser(ctx, &pb.CreateUserRequest{
		Name:  "Rahul",
		Email: "rahul@gmail.com",
	})
	if err != nil {
		log.Fatal("CreateUser error:", err)
	}
	log.Println("Created User:", createRes2.User)

	userID := createRes.User.Id

	getRes, err := client.GetUser(ctx, &pb.GetUserRequest{
		Id: userID,
	})
	if err != nil {
		log.Fatal("GetUser error:", err)
	}
	log.Println("Fetched User:", getRes.User)

	updateRes, err := client.UpdateUser(ctx, &pb.UpdateUserRequest{
		Id:    userID,
		Name:  "Divyam Sinha",
		Email: "divyam.sinha@test.com",
	})
	if err != nil {
		log.Fatal("UpdateUser error:", err)
	}
	log.Println("Updated User:", updateRes.User)
}
