// Package session provides utility functions for managing user sessions in the application.
package session

import (
	"context"
	"errors"

	"github.com/Azzonya/go-shortener/internal/user"
)

// ctxKey is a custom type used as the key type for context values.
type ctxKey int

const (
	// ctxKeyUID is the context key for storing user information.
	ctxKeyUID ctxKey = iota
)

// SetUserContext sets the user information in the context.
func SetUserContext(parent context.Context, u *user.User) context.Context {
	return context.WithValue(parent, ctxKeyUID, u)
}

// GetUserFromContext retrieves the user information from the context.
func GetUserFromContext(ctx context.Context) (u *user.User, ok bool) {
	u, ok = ctx.Value(ctxKeyUID).(*user.User)
	return
}

// GetUser retrieves the user information from the context and performs validation.
func GetUser(c context.Context) (*user.User, error) {
	u, ok := GetUserFromContext(c)
	if !ok {
		return nil, errors.New("failed to get user from context")
	}

	if u.IsNew() {

		return nil, errors.New("no authorized")
	}

	return u, nil
}
