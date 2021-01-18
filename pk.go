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


type Store interface {
	Add(ctx context.Context, account Account) error
	Get(ctx context.Context, name, username string) (account Account, err error)
	List(ctx context.Context) (accounts []Account, err error)
	Delete(ctx context.Context, name, username string) (err error)
	Update(ctx context.Context, account Account) (acc Account, err error)
}
