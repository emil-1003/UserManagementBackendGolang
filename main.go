package main

import (
	"log"

	"github.com/joho/godotenv"

	"github.com/emil-1003/UserManagementBackendGolang/pkg/database"
	"github.com/emil-1003/UserManagementBackendGolang/pkg/server"
)

const (
	apiPath    = "api"
	apiVersion = "v1"
	apiPort    = ":8585"
	apiName    = "UserManagement"
)

func init() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	// Start db connection
	if err := database.ConnectToDb(); err != nil {
		log.Fatalf("Error connecting to database: %s", err)
	}
}

func main() {
	// Start server
	srv, err := server.New(apiName, apiVersion, apiPort, apiPath)
	if err != nil {
		log.Fatalf("Server error: %s", err)
	}

	log.Printf("Starting %s version %s, listening on %s", srv.Name, srv.Version, srv.Port)
	log.Fatal(srv.ListenAndServe())
}
