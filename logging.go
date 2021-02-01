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
	"time"
)

var _ PasswordKeeper = (*loggingMiddleware)(nil)

type loggingMiddleware struct {
	logger *log.Logger
	next   PasswordKeeper
}

func LoggingMiddleware(logger *log.Logger) Middleware {
	return func(keeper PasswordKeeper) PasswordKeeper {
		return &loggingMiddleware{next: keeper, logger: logger}
	}
}

func (l loggingMiddleware) Register(ctx context.Context, username, email, password string) (err error) {
	defer func(begin time.Time) {
		l.logger.Printf("method: register took: %v to register user with id: %v and return err: %v\n",
			time.Since(begin), username, err)
	}(time.Now())

	err = l.next.Register(ctx, username, email, password)
	return
}

func (l loggingMiddleware) Login(ctx context.Context, username, password string) (token string, err error) {
	defer func(begin time.Time) {
		l.logger.Printf("method: login took: %v to generate token for user with id: %v and return err: %v\n",
			time.Since(begin), username, err)
	}(time.Now())

	token, err = l.next.Login(ctx, username, password)
	return
}

func (l loggingMiddleware) Add(ctx context.Context, token string, account Account) (err error) {
	defer func(begin time.Time) {
		l.logger.Printf("method: add took: %v to add new user with id: %v and return err: %v\n",
			time.Since(begin), account.UserName, err)
	}(time.Now())

	err = l.next.Add(ctx, token, account)
	return
}

func (l loggingMiddleware) Get(ctx context.Context, token, name, username string) (account Account, err error) {
	defer func(begin time.Time) {
		l.logger.Printf("method: get took: %v to retrieve user with id: %v \n",
			time.Since(begin), username)
	}(time.Now())

	account, err = l.next.Get(ctx, token, name, username)
	return
}

func (l loggingMiddleware) Delete(ctx context.Context, token, name, username string) (err error) {
	defer func(begin time.Time) {
		l.logger.Printf("method: delete took: %v to delete user with id: %v and returned err : %v\n",
			time.Since(begin), username, err)
	}(time.Now())

	err = l.next.Delete(ctx, token, name, username)
	return
}

func (l loggingMiddleware) List(ctx context.Context, token string, args map[string]interface{}) (accounts []Account, err error) {
	defer func(begin time.Time) {
		l.logger.Printf("method: list took: %v to retrieve all users (%v) and returned with err: %v \n",
			time.Since(begin), len(accounts), err)
	}(time.Now())

	accounts, err = l.next.List(ctx, token, args)
	return
}

func (l loggingMiddleware) Update(ctx context.Context, token, name, username, account Account) (acc Account, err error) {
	defer func(begin time.Time) {
		l.logger.Printf("method: update took: %v to update acc with username: %v \n",
			time.Since(begin), username)
	}(time.Now())

	acc, err = l.next.Update(ctx, token, name, username, account)
	return
}

func (l loggingMiddleware) AddAll(ctx context.Context, token string, accounts []Account) (err error) {
	defer func(begin time.Time) {
		l.logger.Printf("method: addAll() took: %v to add %v users\n",
			time.Since(begin), len(accounts))
	}(time.Now())

	err = l.next.AddAll(ctx, token, accounts)
	return
}

func (l loggingMiddleware) DeleteAll(ctx context.Context, token string, args map[string]interface{}) (err error) {
	defer func(begin time.Time) {
		l.logger.Printf("method: deleteAll() took: %v to delete all users fitting into %v\n",
			time.Since(begin), args)
	}(time.Now())

	err = l.next.DeleteAll(ctx, token, args)
	return
}
