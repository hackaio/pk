package pk

import (
	"context"
)

type authMiddleware struct {
	userName string
	password string
	next     Service
}

func (a authMiddleware) Init(ctx context.Context, username, email, password string) (err error) {
	panic("implement me")
}

func (a authMiddleware) Add(ctx context.Context, account Account) (err error) {
	panic("implement me")
}

func (a authMiddleware) Get(ctx context.Context, username string) (acc Account, err error) {
	panic("implement me")
}

func (a authMiddleware) List(ctx context.Context) (accounts []Account, err error) {
	panic("implement me")
}

func (a authMiddleware) Delete(ctx context.Context, name, username string) (err error) {
	panic("implement me")
}

func (a authMiddleware) Update(ctx context.Context, account Account) (acc Account, err error) {
	panic("implement me")
}

// AuthMiddleware takes a username and password as a dependency
// and returns a Service Middleware.
func AuthMiddleware(username, password string) Middleware {
	return func(next Service) Service {
		return &authMiddleware{
			userName: username,
			password: password,
			next:     next,
		}
	}

}
