package main

import (
	"context"
	"encoding/base64"
	"log"
	"time"

	pb "grpc-crud-proj/proto/userpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
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

	// Helper to add auth
	auth := base64.StdEncoding.EncodeToString([]byte("admin:password"))
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Basic "+auth)

	createRes, err := client.CreateUser(ctx, &pb.CreateUserRequest{
		Name:  "Divyam",
		Email: "divyam@test.com",
	})
	if err != nil {
		log.Fatal("CreateUser error:", err)
	}
	log.Println("Created User:", createRes.User)

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

	log.Println("Testing DeleteUser with 2s deadline (server sleeps 5m)...")
	deleteRes, err := client.DeleteUser(ctx, &pb.DeleteUserRequest{
		Id: userID,
	})
	if err != nil {
		log.Printf("DeleteUser expected error (deadline): %v", err)
	} else {
		log.Println("DeleteUser result:", deleteRes.Message)
	}
}
