package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"google.golang.org/api/idtoken"

	"github.com/nickolasgough/cloud-9-iam/internal/api"
	"github.com/nickolasgough/cloud-9-iam/internal/auth"
	"github.com/nickolasgough/cloud-9-iam/internal/google"
	"github.com/nickolasgough/cloud-9-iam/internal/shared/constants"
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
	mux.HandleFunc("/sign-in/with-password", corsInterceptor(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil {
			fmt.Printf("Error occurred getting session cookie: %s\n", err)
			fmt.Printf("Cookies length is: %d\n", len(r.Cookies()))
			return
		}
		fmt.Printf("Session cookie is: %+v\n", cookie)
	}, []string{http.MethodPost}))

	fmt.Printf("Server listening on port %d\n", PORT)
	err = http.ListenAndServe(fmt.Sprintf(":%d", PORT), mux)
	if err != nil {
		os.Exit(1)
	}
}

func corsInterceptor(handlerFn api.HandlerFn, supportedMethods []string) api.HandlerFn {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(supportedMethods, ", "))
		if r.Header.Get("Access-Control-Request-Headers") != "" {
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		}
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			status := http.StatusOK
			statusText := http.StatusText(http.StatusOK)
			w.WriteHeader(status)
			w.Write([]byte(statusText))
			return
		}

		handlerFn(w, r)
	}
}
