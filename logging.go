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

type loggingMiddleware struct {
	logger *log.Logger
	next   PasswordKeeper
}

func (l loggingMiddleware) AddMany(ctx context.Context, req BulkAddRequest) (err error) {
	defer func(begin time.Time) {
		l.logger.Printf("method: add-many took: %v to add %d users and return err: %v\n",
			time.Since(begin),len(req.Accounts),err)
	}(time.Now())

	err = l.next.AddMany(ctx,req)
	return
}

func (l loggingMiddleware) CredStore() CredStore {
	panic("implement me")
}

func LoggingMiddleware(logger *log.Logger) Middleware {
	return func(keeper PasswordKeeper) PasswordKeeper {
		return &loggingMiddleware{next: keeper}
	}
}

func (l loggingMiddleware) Register(ctx context.Context, request RegisterRequest) (err ErrResponse) {
	defer func(begin time.Time) {
		l.logger.Printf("method: register took: %v to register user with id: %v and return err: %v\n",
			time.Since(begin),request.Username,err.Err)
	}(time.Now())

	err = l.next.Register(ctx,request)
	return
}

func (l loggingMiddleware) Login(ctx context.Context, request LoginRequest) (response LoginResponse) {
	defer func(begin time.Time) {
		l.logger.Printf("method: login took: %v to generate token for user with id: %v and return err: %v\n",
			time.Since(begin),request.UserName,response.Err)
	}(time.Now())

	response = l.next.Login(ctx,request)
	return
}

func (l loggingMiddleware) Add(ctx context.Context, request AddRequest) (err ErrResponse) {
	defer func(begin time.Time) {
		l.logger.Printf("method: add took: %v to add new user with id: %v and return err: %v\n",
			time.Since(begin),request.UserName,err.Err)
	}(time.Now())

	err = l.next.Add(ctx,request)
	return
}

func (l loggingMiddleware) Get(ctx context.Context, request GetRequest) (response GetResponse) {
	defer func(begin time.Time) {
		l.logger.Printf("method: get took: %v to retrieve user with id: %v \n",
			time.Since(begin),request.UserName)
	}(time.Now())

	response = l.next.Get(ctx,request)
	return
}

func (l loggingMiddleware) Delete(ctx context.Context, request GetRequest) (err ErrResponse) {
	defer func(begin time.Time) {
		l.logger.Printf("method: delete took: %v to delete user with id: %v and returned err : %v\n",
			time.Since(begin),request.UserName,err.Err)
	}(time.Now())

	err = l.next.Delete(ctx,request)
	return
}

func (l loggingMiddleware) List(ctx context.Context, req ListRequest) (list ListResponse) {
	defer func(begin time.Time) {
		l.logger.Printf("method: list took: %v to retrieve all users (%v) and returned with err: %v \n",
			time.Since(begin),len(list.Accounts),list.Err)
	}(time.Now())

	list = l.next.List(ctx,req)
	return
}

func (l loggingMiddleware) Update(ctx context.Context, request UpdateRequest) (response ErrResponse) {
	defer func(begin time.Time) {
		l.logger.Printf("method: update took: %v to update acc with username: %v \n",
			time.Since(begin),request.Username)
	}(time.Now())

	response = l.next.Update(ctx,request)
	return
}
