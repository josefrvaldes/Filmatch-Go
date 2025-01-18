package firebase

import (
	"context"
	"filmatch/interfaces"

	"firebase.google.com/go/v4/auth"
)

var _ interfaces.AuthClient = (*FirebaseAuthClient)(nil) // Ensure FirebaseAuthClient implements AuthClient

type FirebaseAuthClient struct { // implements AuthClient
	client *auth.Client
}

func (f *FirebaseAuthClient) VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	return f.client.VerifyIDToken(ctx, idToken)
}

func NewFirebaseAuthClient(client *auth.Client) interfaces.AuthClient {
	return &FirebaseAuthClient{client: client}
}
