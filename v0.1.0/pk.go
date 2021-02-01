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

package v0_1_0

import (
	"context"
	"github.com/hackaio/pk"
)

var _ PasswordKeeper = (*passwordKeeper)(nil)

//Account all details of account
type Account struct {
	Name     string `json:"name,omitempty"`
	UserName string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
	Created  string `json:"created,omitempty"`
}

//DBAccount it how the account is stored in database
type DBAccount struct {
	Name      string `json:"name,omitempty"`
	UserName  string `json:"username,omitempty"`
	Email     string `json:"email,omitempty"`
	Hash      string `json:"hash,omitempty"`
	Encoded   []byte `json:"encoded,omitempty"`
	Digest    []byte `json:"digest,omitempty"`
	Signature []byte `json:"signature,omitempty"`
	Created   string `json:"created,omitempty"`
}

// PasswordKeeper describes the service.
type PasswordKeeper interface {

	//Register creates a new account. The function takes
	//Username, Email, Password. Account name is "master"
	//This is the first function to be called when running
	//the app
	Register(ctx context.Context, username, email, password string) (err error)

	//Login function returns token after a user has supplied
	//his correct username and password or else an error
	//It also works fine if the email is supplied in place of
	//username
	Login(ctx context.Context, username, password string) (token string, err error)

	//Add a new account. It takes token and new account details.
	//The method returns err if the process is not allowed
	Add(ctx context.Context, token string, account Account) (err error)

	//Get returns details of the account after a user has supplied token
	//username and name of the account
	//username e.g pius
	//name e.g github
	Get(ctx context.Context, token, name, username string) (account Account, err error)

	//Delete removes the details of the account after a user has supplied token
	//username and name of the account
	//username e.g pius
	//name e.g github
	Delete(ctx context.Context, token, name, username string) (err error)

	//List returns all the accounts registered under the master
	//accounts
	List(ctx context.Context, token string, args map[string]interface{}) (accounts []Account, err error)

	//Updates the details of the account
	//name and username of the account as of right now
	//Account should have new username and password or email
	Update(ctx context.Context, token, name, username, account Account) (acc Account, err error)

	//AddAll is an API for bulk addition. where a lot of accounts are added all at once
	AddAll(ctx context.Context, token string, accounts []Account) (err error)

	//DeleteAll contains API for bulk Delete
	//e.g You want to delete all accounts by name of instagram
	//or you want to delete all accounts registered under a certain email address
	DeleteAll(ctx context.Context, token string, args map[string]interface{}) (err error)
}

type passwordKeeper struct {
	hasher      pk.Hasher
	passwords   pk.PasswordStore
	tokenizer   pk.Tokenizer
	es          pk.EncoderSigner
}

func NewPasswordKeeper(
	hasher pk.Hasher, store pk.PasswordStore,
	tokenizer pk.Tokenizer, es pk.EncoderSigner) PasswordKeeper {
	return &passwordKeeper{
		hasher:      hasher,
		passwords:   store,
		tokenizer:   tokenizer,
		es:          es,
	}
}

func (p passwordKeeper) Register(ctx context.Context, username, email, password string) (err error) {
	panic("implement me")
}

func (p passwordKeeper) Login(ctx context.Context, username, password string) (token string, err error) {
	panic("implement me")
}

func (p passwordKeeper) Add(ctx context.Context, token string, account Account) (err error) {
	panic("implement me")
}

func (p passwordKeeper) Get(ctx context.Context, token, name, username string) (account Account, err error) {
	panic("implement me")
}

func (p passwordKeeper) Delete(ctx context.Context, token, name, username string) (err error) {
	panic("implement me")
}

func (p passwordKeeper) List(ctx context.Context, token string, args map[string]interface{}) (accounts []Account, err error) {
	panic("implement me")
}

func (p passwordKeeper) Update(ctx context.Context, token, name, username, account Account) (acc Account, err error) {
	panic("implement me")
}

func (p passwordKeeper) AddAll(ctx context.Context, token string, accounts []Account) (err error) {
	panic("implement me")
}

func (p passwordKeeper) DeleteAll(ctx context.Context, token string, args map[string]interface{}) (err error) {
	panic("implement me")
}



type Middleware func(keeper PasswordKeeper) PasswordKeeper

func New( hasher pk.Hasher, store pk.PasswordStore, tokenizer pk.Tokenizer,
	es pk.EncoderSigner, middlewares []Middleware) PasswordKeeper {

	var keeper = NewPasswordKeeper(hasher, store, tokenizer, es)

	for _, middleware := range middlewares {
		keeper = middleware(keeper)
	}

	return keeper
}

func AddMiddlewares(keeper PasswordKeeper, middlewares []Middleware) PasswordKeeper {
	for _, middleware := range middlewares {
		keeper = middleware(keeper)
	}
	return keeper
}



