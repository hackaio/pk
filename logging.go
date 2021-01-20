/*
 * Copyright © 2021 PIUS ALFRED me.pius1102@gmail.com
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *     http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
