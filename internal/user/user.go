// Package user provides functionality for managing user information.
//
// This package includes a User struct representing a user and methods for working with user objects.
//
// Users can be created with or without specifying an ID. If an ID is not provided during creation, a new
// universally unique identifier (UUID) will be generated for the user.
package user

import (
	"encoding/base64"
	"fmt"

	"github.com/gofrs/uuid/v5"
)

// User represents a user with an ID and a flag indicating whether it's a new user.
type User struct {
	ID  string
	new bool
}

// IsNew returns true if the user is newly created, otherwise false.
func (u *User) IsNew() bool {
	return u.new
}

// NewWithID creates a new user with the specified ID.
func NewWithID(id string) *User {
	return &User{ID: id}
}

// New creates a new user with a randomly generated ID.
func New() (*User, error) {
	uid, err := uuid.NewV4()
	if err != nil {
		return nil, fmt.Errorf("cannot generate new uuid: %w", err)
	}
	return &User{
		ID:  base64.RawURLEncoding.EncodeToString(uid.Bytes()),
		new: true,
	}, nil
}
