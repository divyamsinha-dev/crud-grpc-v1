package main

/*
 * ============================================
 * gRPC SERVER WITH HTTP GATEWAY
 * ============================================
 *
 * This file contains the server that:
 * 1. Runs a gRPC server (for gRPC clients)
 * 2. Runs an HTTP gateway (for Postman/browser)
 *
 * Think of it like a restaurant:
 * - gRPC server = Kitchen (where food is made)
 * - HTTP gateway = Waiter (translates orders from customers to kitchen)
 * - Database = Storage (where ingredients/data are kept)
 */

import (
	"context"      // Used for request context (like request timeout, cancellation)
	"database/sql" // For database operations
	"log"          // For logging messages
	"net"          // For network operations (listening on ports)
	"net/http"     // For HTTP server (Postman/browser requests)

	"grpc-crud-proj/db"              // Our database connection package
	gw "grpc-crud-proj/proto/userpb" // Gateway code (same package, different use)
	pb "grpc-crud-proj/proto/userpb" // Generated code from proto file (pb = protobuf)

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime" // HTTP to gRPC translator
	"google.golang.org/grpc"                            // gRPC library
	"google.golang.org/grpc/credentials/insecure"       // For local development (no SSL)
)

/*
 * ============================================
 * SERVER STRUCT
 * ============================================
 *
 * This is like a "class" that holds our server logic.
 * It has:
 * - UnimplementedUserServiceServer: Base implementation (required by gRPC)
 * - db: Database connection to PostgreSQL
 *
 * Think of it as a "service object" that has access to the database.
 */
type server struct {
	pb.UnimplementedUserServiceServer         // Required: Base implementation from generated code
	db                                *sql.DB // Our database connection
}

/*
 * ============================================
 * CREATE USER FUNCTION
 * ============================================
 *
 * This function is called when someone wants to create a new user.
 *
 * Flow:
 * 1. Receive request with name and email
 * 2. Insert into database
 * 3. Get the auto-generated ID
 * 4. Return the created user
 *
 * Real-life analogy: Like filling out a form to register a new account
 */
func (s *server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	// Step 1: Prepare a variable to store the new user's ID
	var id int

	// Step 2: Execute SQL INSERT query
	// $1 and $2 are placeholders for req.Name and req.Email (prevents SQL injection)
	// RETURNING id means: "Give me back the ID that was auto-generated"
	err := s.db.QueryRow(
		"INSERT INTO users(name,email) VALUES($1,$2) RETURNING id",
		req.Name,  // $1 = name from request
		req.Email, // $2 = email from request
	).Scan(&id) // Scan stores the returned ID into our 'id' variable

	// Step 3: Check if there was an error (like duplicate email)
	if err != nil {
		return nil, err // Return error to caller
	}

	// Step 4: Create and return the response
	// &pb.UserResponse means: create a pointer to UserResponse struct
	return &pb.UserResponse{
		User: &pb.User{ // Create a User object inside the response
			Id:    int32(id), // Convert int to int32 (database returns int, proto expects int32)
			Name:  req.Name,  // Use the name from request
			Email: req.Email, // Use the email from request
		},
	}, nil // nil means no error
}

/*
 * ============================================
 * GET USER FUNCTION
 * ============================================
 *
 * This function retrieves a user by their ID.
 *
 * Flow:
 * 1. Receive request with user ID
 * 2. Query database for that user
 * 3. Return user data
 *
 * Real-life analogy: Like looking up someone's profile by their ID number
 */
func (s *server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	// Step 1: Create an empty User object to store the result
	var user pb.User

	// Step 2: Execute SQL SELECT query
	// $1 is placeholder for req.Id
	err := s.db.QueryRow(
		"SELECT id, name, email FROM users WHERE id=$1",
		req.Id, // $1 = ID from request
	).Scan(&user.Id, &user.Name, &user.Email) // Scan fills our user object with database results

	// Step 3: Check if user was found (if not, err will be "no rows")
	if err != nil {
		return nil, err // Return error (like "user not found")
	}

	// Step 4: Return the user wrapped in a UserResponse
	return &pb.UserResponse{User: &user}, nil
}

/*
 * ============================================
 * UPDATE USER FUNCTION
 * ============================================
 *
 * This function updates an existing user's information.
 *
 * Flow:
 * 1. Receive request with ID, new name, and new email
 * 2. Update database record
 * 3. Return updated user data
 *
 * Real-life analogy: Like updating your profile information
 */
func (s *server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	// Step 1: Execute SQL UPDATE query
	// $1 = new name, $2 = new email, $3 = user ID
	// _ means we ignore the result (we don't need it)
	_, err := s.db.Exec(
		"UPDATE users SET name=$1, email=$2 WHERE id=$3",
		req.Name,  // $1 = new name
		req.Email, // $2 = new email
		req.Id,    // $3 = which user to update
	)

	// Step 2: Check for errors
	if err != nil {
		return nil, err
	}

	// Step 3: Return the updated user (we create it from the request data)
	return &pb.UserResponse{
		User: &pb.User{
			Id:    req.Id,    // Same ID
			Name:  req.Name,  // New name
			Email: req.Email, // New email
		},
	}, nil
}

/*
 * ============================================
 * DELETE USER FUNCTION
 * ============================================
 *
 * This function deletes a user from the database.
 *
 * Flow:
 * 1. Receive request with user ID
 * 2. Delete from database
 * 3. Return success message
 *
 * Real-life analogy: Like deleting an account
 */
func (s *server) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	// Step 1: Execute SQL DELETE query
	// $1 = user ID to delete
	_, err := s.db.Exec("DELETE FROM users WHERE id=$1", req.Id)

	// Step 2: Check for errors
	if err != nil {
		return nil, err
	}

	// Step 3: Return success message
	return &pb.DeleteUserResponse{
		Message: "User deleted",
	}, nil
}

/*
 * ============================================
 * MAIN FUNCTION - SERVER STARTUP
 * ============================================
 *
 * This is where everything starts. Think of it as the "power button" for the server.
 *
 * It does two things:
 * 1. Starts gRPC server (for gRPC clients)
 * 2. Starts HTTP gateway (for Postman/browser)
 *
 * Real-life analogy: Like opening a restaurant - you need both the kitchen (gRPC)
 * and the front door (HTTP gateway) to be open.
 */
func main() {
	// ============================================
	// STEP 1: Connect to Database
	// ============================================
	// This opens a connection to PostgreSQL
	// Think of it as connecting to a storage warehouse
	dbConn := db.Connect()

	// ============================================
	// STEP 2: Start gRPC Server (in background)
	// ============================================
	// We use 'go func()' to run this in a separate "goroutine" (like a separate thread)
	// This allows both gRPC and HTTP to run at the same time
	go func() {
		// 2a. Listen on port 50051 for gRPC connections
		// Think of this as opening a door on port 50051
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatal("Failed to listen on gRPC port:", err)
		}

		// 2b. Create a new gRPC server
		// This is like creating a kitchen that can handle gRPC orders
		grpcServer := grpc.NewServer()

		// 2c. Register our service with the gRPC server
		// This tells gRPC: "Hey, when someone calls UserService methods, use our server struct"
		// Think of it as telling the kitchen: "Here's the menu and recipes"
		pb.RegisterUserServiceServer(grpcServer, &server{db: dbConn})

		// 2d. Start serving (this blocks, so it runs forever)
		// This is like the kitchen starting to work - it keeps running
		log.Println("🚀 gRPC server running on :50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal("Failed to serve gRPC:", err)
		}
	}()

	// ============================================
	// STEP 3: Setup HTTP Gateway (for Postman)
	// ============================================
	// The gateway translates HTTP/JSON requests to gRPC calls
	// Think of it as a translator between customers (Postman) and kitchen (gRPC server)

	// 3a. Create a context (used for request handling)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel() // Clean up when function ends

	// 3b. Connect to our own gRPC server
	// This is interesting: The gateway connects to the gRPC server we just started
	// Think of it as the waiter connecting to the kitchen intercom
	conn, err := grpc.NewClient(
		"localhost:50051",                                        // Address of our gRPC server
		grpc.WithTransportCredentials(insecure.NewCredentials()), // No SSL for local dev
	)
	if err != nil {
		log.Fatal("Failed to dial gRPC server:", err)
	}
	defer conn.Close() // Close connection when done

	// 3c. Create HTTP router (this handles HTTP requests)
	// Think of it as the waiter who takes orders
	mux := runtime.NewServeMux()

	// 3d. Register our service handlers with the gateway
	// This tells the gateway: "When someone calls /v1/users, translate it to CreateUser RPC"
	// Think of it as teaching the waiter: "When customer says 'POST /v1/users',
	// translate that to 'CreateUser' order for the kitchen"
	err = gw.RegisterUserServiceHandler(ctx, mux, conn)
	if err != nil {
		log.Fatal("Failed to register gateway:", err)
	}

	// ============================================
	// STEP 4: Start HTTP Server
	// ============================================
	// This starts listening on port 8080 for HTTP requests (from Postman/browser)
	// Think of it as opening the restaurant's front door
	log.Println("🌐 HTTP/REST gateway running on :8080")
	log.Println("📝 You can now use Postman to test:")
	log.Println("   POST   http://localhost:8080/v1/users")
	log.Println("   GET    http://localhost:8080/v1/users/{id}")
	log.Println("   PUT    http://localhost:8080/v1/users/{id}")
	log.Println("   DELETE http://localhost:8080/v1/users/{id}")

	// Start the HTTP server (this blocks, so it runs forever)
	// This is like the restaurant staying open - it keeps accepting customers
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal("Failed to serve HTTP:", err)
	}

	// Note: We never reach here because ListenAndServe runs forever
	// The program only exits if there's an error
}
