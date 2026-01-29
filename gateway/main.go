package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	pb "grpc-crud-proj/proto/userpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

/*
 * 🎓 CONCEPT: API Gateway Pattern
 * 
 * Real-life Example: Think of a restaurant with two types of customers:
 * 1. Regular customers (browsers) - speak English (HTTP/JSON)
 * 2. VIP customers (microservices) - speak French (gRPC)
 * 
 * The API Gateway is like a bilingual waiter who:
 * - Takes orders in English from regular customers (HTTP requests)
 * - Translates them to French for the kitchen staff (gRPC calls)
 * - Translates responses back to English for the customer
 * 
 * This way, the kitchen (gRPC services) stays efficient, but everyone can order!
 */

// GatewayServer holds the gRPC client connection
type GatewayServer struct {
	grpcClient pb.UserServiceClient
	conn       *grpc.ClientConn
}

// NewGatewayServer creates a new gateway server with gRPC connection
func NewGatewayServer() (*GatewayServer, error) {
	// Connect to the gRPC server (like connecting to the kitchen)
	conn, err := grpc.Dial(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	client := pb.NewUserServiceClient(conn)

	return &GatewayServer{
		grpcClient: client,
		conn:       conn,
	}, nil
}

// Close closes the gRPC connection
func (g *GatewayServer) Close() {
	g.conn.Close()
}

/*
 * 🎓 CONCEPT: CORS (Cross-Origin Resource Sharing)
 * 
 * Real-life Example: Imagine a library (your API) that only allows
 * people from your city to check out books. CORS is like the librarian
 * who checks your ID and decides if you're allowed.
 * 
 * Browsers enforce CORS - they won't let JavaScript from one website
 * (origin) make requests to another website unless the server explicitly
 * allows it. This is a security feature!
 */
func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

// handleOptions handles preflight CORS requests
func (g *GatewayServer) handleOptions(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	w.WriteHeader(http.StatusOK)
}

/*
 * 🎓 CONCEPT: HTTP Request/Response Cycle
 * 
 * Real-life Example: Ordering food delivery
 * 1. You (browser) make a request: "I want pizza" (HTTP POST)
 * 2. The restaurant (server) processes it: "Got it, making pizza"
 * 3. The restaurant responds: "Here's your pizza" (HTTP 200 + JSON)
 * 
 * HTTP Methods:
 * - GET: "Show me something" (like viewing a menu)
 * - POST: "Create something new" (like placing an order)
 * - PUT: "Update something" (like changing your order)
 * - DELETE: "Remove something" (like canceling an order)
 */

// CreateUserHandler handles HTTP POST /api/users
func (g *GatewayServer) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse JSON from HTTP request body
	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Create context with timeout (like setting a timer for the order)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 🎓 CONCEPT: HTTP to gRPC Translation
	// Convert HTTP request → gRPC request → gRPC response → HTTP response
	grpcReq := &pb.CreateUserRequest{
		Name:  req.Name,
		Email: req.Email,
	}

	// Call gRPC service (like sending order to kitchen)
	grpcResp, err := g.grpcClient.CreateUser(ctx, grpcReq)
	if err != nil {
		http.Error(w, "gRPC error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert gRPC response to JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(grpcResp.User)
}

// GetUserHandler handles HTTP GET /api/users/:id
func (g *GatewayServer) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from URL path (like reading the order number)
	idStr := r.URL.Path[len("/api/users/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	grpcReq := &pb.GetUserRequest{
		Id: int32(id),
	}

	grpcResp, err := g.grpcClient.GetUser(ctx, grpcReq)
	if err != nil {
		http.Error(w, "gRPC error: "+err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(grpcResp.User)
}

// UpdateUserHandler handles HTTP PUT /api/users/:id
func (g *GatewayServer) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from URL
	idStr := r.URL.Path[len("/api/users/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Parse JSON body
	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	grpcReq := &pb.UpdateUserRequest{
		Id:    int32(id),
		Name:  req.Name,
		Email: req.Email,
	}

	grpcResp, err := g.grpcClient.UpdateUser(ctx, grpcReq)
	if err != nil {
		http.Error(w, "gRPC error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(grpcResp.User)
}

// DeleteUserHandler handles HTTP DELETE /api/users/:id
func (g *GatewayServer) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Path[len("/api/users/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	grpcReq := &pb.DeleteUserRequest{
		Id: int32(id),
	}

	grpcResp, err := g.grpcClient.DeleteUser(ctx, grpcReq)
	if err != nil {
		http.Error(w, "gRPC error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": grpcResp.Message})
}

/*
 * 🎓 CONCEPT: HTTP Routing
 * 
 * Real-life Example: A receptionist at a hotel who directs guests:
 * - "Room 101? Go to floor 1" (route /api/users/1 → GetUserHandler)
 * - "Check-in? Go to front desk" (route /api/users → CreateUserHandler)
 * 
 * We use a simple pattern matching to route requests to the right handler.
 */
func (g *GatewayServer) setupRoutes() {
	// Serve static files (HTML, CSS, JS) from the static directory
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	// API routes
	http.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			g.handleOptions(w, r)
			return
		}
		if r.Method == http.MethodPost {
			g.CreateUserHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Route for /api/users/:id (GET, PUT, DELETE)
	http.HandleFunc("/api/users/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			g.handleOptions(w, r)
			return
		}
		switch r.Method {
		case http.MethodGet:
			g.GetUserHandler(w, r)
		case http.MethodPut:
			g.UpdateUserHandler(w, r)
		case http.MethodDelete:
			g.DeleteUserHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func main() {
	// Initialize gateway server
	gateway, err := NewGatewayServer()
	if err != nil {
		log.Fatal("Failed to create gateway:", err)
	}
	defer gateway.Close()

	// Setup routes
	gateway.setupRoutes()

	// Start HTTP server on port 8080
	log.Println("🚀 REST API Gateway running on http://localhost:8080")
	log.Println("📱 Open http://localhost:8080 in your browser!")
	log.Println("🔌 Make sure gRPC server is running on :50051")
	
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
