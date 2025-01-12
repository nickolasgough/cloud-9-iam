package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"net/url"
	"os"

	_ "github.com/lib/pq"
	"google.golang.org/api/idtoken"

	signinapi "github.com/nickolasgough/cloud-9-iam/internal/api/sign-in"
	userapi "github.com/nickolasgough/cloud-9-iam/internal/api/user"
	"github.com/nickolasgough/cloud-9-iam/internal/auth"
	"github.com/nickolasgough/cloud-9-iam/internal/google"
	"github.com/nickolasgough/cloud-9-iam/internal/shared/api"
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
	psqlInfo := fmt.Sprintf("postgresql://localhost:5432/cloud-9?user=%s&password=%s&sslmode=disable", url.QueryEscape(psqlUsername), url.QueryEscape(psqlPassword))
	fmt.Printf("Info: %s\n", psqlInfo)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Printf("Failed to open database with: %s\n", err)
		os.Exit(1)
	}
	err = db.Ping()
	if err != nil {
		fmt.Printf("Failed to connect to database with: %s\n", err)
		os.Exit(1)
	}
	defer db.Close()

	// Register internal models
	err = usermodel.Register(db)
	if err != nil {
		fmt.Printf("Failed to register the user model with: %s\n", err)
		os.Exit(1)
	}

	// Register repositories and services
	userRepository := userrepository.NewRepository(db)
	userService := userservice.NewService(userRepository)

	// Register public API endpoints
	mux.HandleFunc("/sign-in/with-google", signinapi.SignInWithGoogle(ctx, gcpClientSecret, googleService, authService))
	mux.HandleFunc("/sign-in/with-password", signinapi.SignInWithPassword(ctx, userService, authService))

	// Register private API endpoints
	handlers := []api.ApiHandler{
		userapi.NewCreateUserHandler(userService),
	}
	for _, h := range handlers {
		mux.HandleFunc(h.Path(), api.ServeHTTP(authService, h))
	}

	fmt.Printf("Server listening on port %d\n", PORT)
	err = http.ListenAndServe(fmt.Sprintf(":%d", PORT), mux)
	if err != nil {
		os.Exit(1)
	}
}
