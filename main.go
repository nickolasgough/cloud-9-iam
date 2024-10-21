package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"google.golang.org/api/idtoken"

	"github.com/nickolasgough/cloud-9-iam/internal/api"
	"github.com/nickolasgough/cloud-9-iam/internal/auth"
	"github.com/nickolasgough/cloud-9-iam/internal/google"
	"github.com/nickolasgough/cloud-9-iam/internal/shared/constants"
	usermodel "github.com/nickolasgough/cloud-9-iam/internal/user/model"
	userrepository "github.com/nickolasgough/cloud-9-iam/internal/user/repository"
	userservice "github.com/nickolasgough/cloud-9-iam/internal/user/service"
)

const (
	PORT = 8000
)

func main() {
	ctx := context.Background()
	mux := http.NewServeMux()

	// Initialize environment data.
	gcpClientID := os.Getenv(constants.GCP_CLIENT_ID)
	gcpClientSecret := os.Getenv(constants.GCP_CLIENT_SECRET)
	if gcpClientID == "" || gcpClientSecret == "" {
		fmt.Printf("Failed to load GCP_CLIENT_ID and/or GCP_CLIENT_SECRET environment variable\n")
		os.Exit(1)
	}
	psqlUsername := os.Getenv(constants.PSQL_USERNAME)
	psqlPassword := os.Getenv(constants.PSQL_PASSWORD)
	if gcpClientID == "" || gcpClientSecret == "" {
		fmt.Printf("Failed to load PSQL_USERNAME and/or PSQL_PASSWORD environment variable\n")
		os.Exit(1)
	}

	// Initialize the Google ID token validator client
	googleIDVerifier, err := idtoken.NewValidator(ctx)
	if err != nil {
		fmt.Printf("Failed to initialize Google ID validator client with: %s\n", err)
		os.Exit(1)
	}
	googleService := google.NewService(gcpClientID, googleIDVerifier)

	// Initialize the Authentication service
	authService := auth.NewService(gcpClientSecret)

	// Initialize the database connection
	dsn := fmt.Sprintf("host=localhost user=%s password=%s dbname=cloud-9 port=5432 sslmode=disable", psqlUsername, psqlPassword)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		fmt.Printf("Failed to open database with: %s\n", err)
		os.Exit(1)
	}
	err = db.Ping()
	if err != nil {
		fmt.Printf("Failed to connect to database with: %s\n", err)
		os.Exit(1)
	}

	// Register internal models
	err = usermodel.Register(db)
	if err != nil {
		fmt.Printf("Failed to register the user model with: %s\n", err)
		os.Exit(1)
	}

	// Register repositories and services
	userRepository := userrepository.NewRepository(db)
	userService := userservice.NewService(userRepository)

	// Register internal API endpoints
	handlers := []*api.ApiEndpoint{
		api.NewApiEndpoint(api.NewCreateUserHandler(userService)),
	}
	for _, handler := range handlers {
		mux.Handle(handler.Path(), handler)
	}

	// Register external API endpoints
	mux.HandleFunc("/sign-in/with-google", api.SignInWithGoogle(ctx, gcpClientSecret, googleService, authService))

	fmt.Printf("Server listening on port %d\n", PORT)
	err = http.ListenAndServe(fmt.Sprintf(":%d", PORT), mux)
	if err != nil {
		os.Exit(1)
	}
}
