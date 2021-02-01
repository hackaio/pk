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

var _ PasswordKeeper = (*loggingMiddleware)(nil)

type loggingMiddleware struct {
	logger *log.Logger
	next   PasswordKeeper
}

func LoggingMiddleware(logger *log.Logger) Middleware {
	return func(keeper PasswordKeeper) PasswordKeeper {
		return &loggingMiddleware{next: keeper,logger: logger}
	}
}

func (l loggingMiddleware) Register(ctx context.Context, username, email, password string) (err error) {
	panic("implement me")
}

func (l loggingMiddleware) Login(ctx context.Context, username, password string) (token string, err error) {
	panic("implement me")
}

func (l loggingMiddleware) Add(ctx context.Context, token string, account Account) (err error) {
	panic("implement me")
}

func (l loggingMiddleware) Get(ctx context.Context, token, name, username string) (account Account, err error) {
	panic("implement me")
}

func (l loggingMiddleware) Delete(ctx context.Context, token, name, username string) (err error) {
	panic("implement me")
}

func (l loggingMiddleware) List(ctx context.Context, token string, args map[string]interface{}) (accounts []Account, err error) {
	panic("implement me")
}

func (l loggingMiddleware) Update(ctx context.Context, token, name, username, account Account) (acc Account, err error) {
	panic("implement me")
}

func (l loggingMiddleware) AddAll(ctx context.Context, token string, accounts []Account) (err error) {
	panic("implement me")
}

func (l loggingMiddleware) DeleteAll(ctx context.Context, token string, args map[string]interface{}) (err error) {
	panic("implement me")
}


