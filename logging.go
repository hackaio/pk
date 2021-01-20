package pk

import (
	"context"
	"log"
)

type loggingMiddleware struct {
	logger *log.Logger
	next   PasswordKeeper
}

func LoggingMiddleware(logger *log.Logger) Middleware {
	return func(keeper PasswordKeeper) PasswordKeeper {
		return &loggingMiddleware{next: keeper}
	}
}

func (l loggingMiddleware) Register(ctx context.Context, request RegisterRequest) (err ErrResponse) {

	panic("implement me")
}

func (l loggingMiddleware) Login(ctx context.Context, request LoginRequest) (response LoginResponse) {
	panic("implement me")
}

func (l loggingMiddleware) Add(ctx context.Context, request AddRequest) (err ErrResponse) {
	panic("implement me")
}

func (l loggingMiddleware) Get(ctx context.Context, request GetRequest) (response GetResponse) {
	panic("implement me")
}

func (l loggingMiddleware) Delete(ctx context.Context, request GetRequest) (err ErrResponse) {
	panic("implement me")
}

func (l loggingMiddleware) List(ctx context.Context) (list ListResponse) {
	panic("implement me")
}

func (l loggingMiddleware) Update(ctx context.Context, request UpdateRequest) (response ErrResponse) {
	panic("implement me")
}
