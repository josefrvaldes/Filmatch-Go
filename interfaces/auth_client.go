package interfaces

import (
	"context"

	"firebase.google.com/go/v4/auth"
)

type AuthClient interface {
	VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error)
}
