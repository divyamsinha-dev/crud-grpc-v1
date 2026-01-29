package db

/*
 * ============================================
 * DATABASE CONNECTION
 * ============================================
 * 
 * This file handles connecting to PostgreSQL database.
 * 
 * Think of it as the "key" to unlock the storage room.
 * Once connected, we can store and retrieve data.
 */

import (
	"database/sql"  // Go's standard database interface
	"log"           // For logging messages
	"os"            // For reading environment variables

	_ "github.com/lib/pq" // PostgreSQL driver (the _ means we import it for side effects only)
)

/*
 * ============================================
 * CONNECT FUNCTION
 * ============================================
 * 
 * This function connects to PostgreSQL database.
 * 
 * Steps:
 * 1. Get connection string (from environment or use default)
 * 2. Open connection
 * 3. Test connection (ping)
 * 4. Return connection object
 * 
 * Real-life analogy: Like connecting to WiFi
 * - Connection string = WiFi password
 * - Ping = Testing if connection works
 */
func Connect() *sql.DB {
	// Step 1: Get connection string from environment variable
	// If DB_URL is set, use it; otherwise use default
	// This allows flexibility (different databases for dev/prod)
	connStr := os.Getenv("DB_URL")
	if connStr == "" {
		// Default connection string for local development
		// Format: postgres://username@host:port/database?options
		connStr = "postgres://divyam.sinha@localhost:5432/postgres?sslmode=disable"
		
		// Alternative connection string (commented out)
		// Use this if you have a password-protected database
		// connStr = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	}

	// Step 2: Open connection to database
	// sql.Open doesn't actually connect yet, it just prepares
	// Think of it as dialing a phone number (but not answering yet)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err) // If we can't even prepare, exit program
	}

	// Step 3: Test the connection (actually connect)
	// Ping sends a small request to check if database is reachable
	// Think of it as saying "Hello?" to see if someone answers
	if err := db.Ping(); err != nil {
		log.Fatal(err) // If ping fails, database is not available
	}

	// Step 4: Connection successful!
	log.Println("Connected to Postgres")
	
	// Return the connection object
	// This can now be used to execute SQL queries
	return db
}

/*
 * ============================================
 * HOW TO USE THIS CONNECTION
 * ============================================
 * 
 * In server/main.go, we call:
 *   dbConn := db.Connect()
 * 
 * Then we can use it like:
 *   dbConn.QueryRow("SELECT ...")
 *   dbConn.Exec("INSERT ...")
 * 
 * Think of it as:
 * - Connect() = Getting the key to storage room
 * - dbConn = The key itself
 * - Using dbConn = Using the key to access storage
 */
