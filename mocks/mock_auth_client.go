package mocks

import (
	"context"
	"errors"
	"filmatch/interfaces"

	"firebase.google.com/go/v4/auth"
)

var _ interfaces.AuthClient = (*MockAuthClient)(nil) // Ensure FirebaseAuthClient implements AuthClient

type MockAuthClient struct{} // implements AuthClient

func (m *MockAuthClient) VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	if idToken == "valid-token" {
		return &auth.Token{
			UID: "test-uid",
			Claims: map[string]interface{}{
				"email": "test@example.com",
			},
		}, nil
	}
	return nil, errors.New("invalid token")
}
