package google

import (
	"context"

	"google.golang.org/api/idtoken"
)

type GoogleIDVerifier interface {
	Validate(ctx context.Context, googleID string, audience string) (*idtoken.Payload, error)
}
