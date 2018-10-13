// Package backend defines an IMAP server backend interface.
package backend

import (
	"context"
	"errors"
)

// ErrInvalidCredentials is returned by Backend.Login when a username or a
// password is incorrect.
var ErrInvalidCredentials = errors.New("Invalid credentials")

// Backend is an IMAP server backend. A backend operation always deals with
// users.
type Backend interface {
	// Login authenticates a user. If the username or the password is incorrect,
	// it returns ErrInvalidCredentials.
	Login(ctx context.Context, remoteAddr, username, password string) (User, error)
}
