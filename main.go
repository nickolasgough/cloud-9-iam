package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"google.golang.org/api/idtoken"

	"github.com/nickolasgough/cloud-community-iam/internal/api"
	"github.com/nickolasgough/cloud-community-iam/internal/auth"
	"github.com/nickolasgough/cloud-community-iam/internal/google"
	"github.com/nickolasgough/cloud-community-iam/internal/shared/constants"
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

	// Initialize the Google ID token validator client
	googleIDVerifier, err := idtoken.NewValidator(ctx)
	if err != nil {
		fmt.Printf("Failed to initialize Google ID validator client")
		os.Exit(1)
	}
	googleService := google.NewService(gcpClientID, googleIDVerifier)

	// Initialize the Authentication service
	authService := auth.NewService(gcpClientSecret)

	// Register API endpoints.
	mux.HandleFunc("/sign-in/with-google", api.SignInWithGoogle(ctx, gcpClientSecret, googleService, authService))

	fmt.Printf("Server listening on port %d\n", PORT)
	err = http.ListenAndServe(fmt.Sprintf(":%d", PORT), mux)
	if err != nil {
		os.Exit(1)
	}
}
