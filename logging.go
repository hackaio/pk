package pk

import (
	"context"
	"log"
)

type Middleware func(service Service) Service

type loggingMiddleware struct {
	logger *log.Logger
	next   Service
}

func LoggingMiddleware(logger *log.Logger) Middleware {
	return func(service Service) Service {
		return &loggingMiddleware{next: service}
	}
}

func (l loggingMiddleware) Init(ctx context.Context, username, email, password string) error {
	panic("implement me")
}

func (l loggingMiddleware) Auth(ctx context.Context, username, password string) (err error) {
	panic("implement me")
}

func (l loggingMiddleware) Add(ctx context.Context, account Account) (err error) {
	panic("implement me")
}

func (l loggingMiddleware) Get(ctx context.Context, name, username string) (account Account, err error) {
	panic("implement me")
}

func (l loggingMiddleware) List(ctx context.Context) (accounts []Account, err error) {
	panic("implement me")
}

func (l loggingMiddleware) Delete(ctx context.Context, name, username string) (err error) {
	panic("implement me")
}

func (l loggingMiddleware) Update(ctx context.Context, account Account) (acc Account, err error) {
	panic("implement me")
}
