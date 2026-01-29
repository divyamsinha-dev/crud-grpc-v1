package main

/*
 * ============================================
 * gRPC CLIENT EXAMPLE
 * ============================================
 * 
 * This file shows how to use the gRPC service from a client.
 * 
 * Think of it as a customer placing orders at a restaurant:
 * - Client = Customer
 * - gRPC Server = Restaurant
 * - Each function call = Placing an order
 * 
 * This demonstrates all CRUD operations:
 * - CREATE: Add a new user
 * - READ: Get a user by ID
 * - UPDATE: Modify a user
 * - DELETE: Remove a user
 */

import (
	"context"  // For request context (timeout, cancellation)
	"log"      // For logging
	"time"     // For timeout duration

	pb "grpc-crud-proj/proto/userpb" // Generated code (pb = protobuf)

	"google.golang.org/grpc"                            // gRPC library
	"google.golang.org/grpc/credentials/insecure"     // For local dev (no SSL)
)

func main() {
	// ============================================
	// STEP 1: Connect to gRPC Server
	// ============================================
	// This is like dialing a phone number to call the restaurant
	// We're connecting to the server running on localhost:50051
	
	conn, err := grpc.Dial(
		"localhost:50051",                          // Server address
		grpc.WithTransportCredentials(insecure.NewCredentials()), // No SSL for local
	)
	if err != nil {
		log.Fatal("failed to connect:", err) // If connection fails, exit program
	}
	defer conn.Close() // Close connection when function ends (cleanup)

	// ============================================
	// STEP 2: Create a Client
	// ============================================
	// This creates a client object that we can use to call server methods
	// Think of it as getting a menu and order form from the restaurant
	
	client := pb.NewUserServiceClient(conn)

	// ============================================
	// STEP 3: Create a Context with Timeout
	// ============================================
	// Context is like a timer for our request
	// If the server takes more than 5 seconds, we'll cancel the request
	// This prevents waiting forever if something goes wrong
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // Cancel the context when done (cleanup)

	// ============================================
	// STEP 4: CREATE USER (C in CRUD)
	// ============================================
	// This creates a new user in the database
	// Think of it as registering a new account
	
	log.Println("📝 Creating user...")
	createRes, err := client.CreateUser(ctx, &pb.CreateUserRequest{
		Name:  "Divyam",           // User's name
		Email: "divyam@test.com",  // User's email
	})
	if err != nil {
		log.Fatal("CreateUser error:", err) // If error, exit
	}
	log.Println("✅ Created User:", createRes.User)
	log.Printf("   ID: %d, Name: %s, Email: %s\n", 
		createRes.User.Id, createRes.User.Name, createRes.User.Email)
	
	// Save the user ID for later operations
	userID := createRes.User.Id

	// Create another user to show it works multiple times
	log.Println("\n📝 Creating another user...")
	createRes2, err := client.CreateUser(ctx, &pb.CreateUserRequest{
		Name:  "Rahul",
		Email: "rahul@gmail.com",
	})
	if err != nil {
		log.Fatal("CreateUser error:", err)
	}
	log.Println("✅ Created User:", createRes2.User)

	// ============================================
	// STEP 5: GET USER (R in CRUD - Read)
	// ============================================
	// This retrieves a user by their ID
	// Think of it as looking up someone's profile
	
	log.Println("\n🔍 Getting user by ID...")
	getRes, err := client.GetUser(ctx, &pb.GetUserRequest{
		Id: userID, // The ID of the user we just created
	})
	if err != nil {
		log.Fatal("GetUser error:", err)
	}
	log.Println("✅ Fetched User:", getRes.User)
	log.Printf("   ID: %d, Name: %s, Email: %s\n", 
		getRes.User.Id, getRes.User.Name, getRes.User.Email)

	// ============================================
	// STEP 6: UPDATE USER (U in CRUD)
	// ============================================
	// This updates an existing user's information
	// Think of it as editing your profile
	
	log.Println("\n✏️  Updating user...")
	updateRes, err := client.UpdateUser(ctx, &pb.UpdateUserRequest{
		Id:    userID,                    // Which user to update
		Name:  "Divyam Sinha",            // New name
		Email: "divyam.sinha@test.com",  // New email
	})
	if err != nil {
		log.Fatal("UpdateUser error:", err)
	}
	log.Println("✅ Updated User:", updateRes.User)
	log.Printf("   ID: %d, Name: %s, Email: %s\n", 
		updateRes.User.Id, updateRes.User.Name, updateRes.User.Email)

	// ============================================
	// STEP 7: DELETE USER (D in CRUD)
	// ============================================
	// This deletes a user from the database
	// Think of it as deleting an account
	// 
	// NOTE: This is commented out so we don't delete the user we just created
	// Uncomment to test deletion
	
	/*
	log.Println("\n🗑️  Deleting user...")
	_, err = client.DeleteUser(ctx, &pb.DeleteUserRequest{
		Id: userID, // Which user to delete
	})
	if err != nil {
		log.Fatal("DeleteUser error:", err)
	}
	log.Println("✅ User deleted successfully")
	*/

	log.Println("\n🎉 All operations completed successfully!")
}

/*
 * ============================================
 * SUMMARY OF WHAT WE DID
 * ============================================
 * 
 * 1. Connected to gRPC server (like calling the restaurant)
 * 2. Created a client (got the menu)
 * 3. Created users (registered accounts)
 * 4. Retrieved a user (looked up profile)
 * 5. Updated a user (edited profile)
 * 6. (Optional) Deleted a user (removed account)
 * 
 * Each operation follows the same pattern:
 * 1. Call client.MethodName(ctx, request)
 * 2. Check for errors
 * 3. Use the response
 * 
 * This is the beauty of gRPC - it's like calling local functions,
 * but they actually run on a remote server!
 */
