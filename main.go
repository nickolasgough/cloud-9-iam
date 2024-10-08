package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"google.golang.org/api/idtoken"
)

const (
	PORT = 8000

	GCP_OAUTH_CLIENT_ID_ENV_VAR = "GCP_OAUTH_CLIENT_ID"
)

func main() {
	ctx := context.Background()
	// Initialize the multiplexor.
	mux := http.NewServeMux()

	// Initialize environment data.
	GCP_OAUTH_CLIENT_ID := os.Getenv(GCP_OAUTH_CLIENT_ID_ENV_VAR)
	if GCP_OAUTH_CLIENT_ID == "" {
		fmt.Printf("Failed to load GCP_OAUTH_CLIENT_ID\n")
		os.Exit(1)
	}

	// Initialize the Google ID token validator client
	googleJWTValidator, err := idtoken.NewValidator(ctx)
	if err != nil {
		fmt.Printf("Failed to initialize Google ID validator client")
		os.Exit(1)
	}

	mux.HandleFunc("/sign-in/with-google", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			fmt.Printf("Error parsing the form: %s\n", err)
			return
		}
		fmt.Printf("Form contains: %+v\n", r.Form)

		googleJWT := r.Form.Get("credential")
		payload, err := googleJWTValidator.Validate(ctx, googleJWT, GCP_OAUTH_CLIENT_ID)
		if err != nil {
			fmt.Printf("Failed to validate Google JWT with: %s\n", err)
			return
		}
		fmt.Printf("Google JWT validation payload: %+v\n", payload)
	})

	fmt.Printf("Server listening on port %d\n", PORT)
	err = http.ListenAndServe(fmt.Sprintf(":%d", PORT), mux)
	if err != nil {
		os.Exit(1)
	}
}
