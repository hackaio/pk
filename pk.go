package pk

import (
	"context"
	"github.com/hackaio/pk/pkg/errors"
)

var (
	ErrNotFound          = errors.New("not found")
	ErrCouldNotCreateAcc = errors.New("could not create account")
	ErrPermissionDenied  = errors.New("permission denied")
)

// Hasher specifies an API for generating hashes of an arbitrary textual
// content.
type Hasher interface {
	// Hash generates the hashed string from plain-text.
	Hash(string) (string, error)

	// Compare compares plain-text version to the hashed one.
	//An error should indicate failed comparison.
	Compare(string, string) error
}

type Account struct {
	Name     string `json:"name,omitempty"`
	UserName string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
	Created  string `json:"created,omitempty"`
}

type Store interface {
	Add(ctx context.Context, account Account) error
	Get(ctx context.Context, name, username string) (account Account, err error)
	List(ctx context.Context) (accounts []Account, err error)
	Delete(ctx context.Context, name, username string) (err error)
	Update(ctx context.Context, account Account) (acc Account, err error)
}
