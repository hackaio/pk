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
	"fmt"
	"github.com/hackaio/pk/pkg/errors"
	"time"
)

const (
	AppName = "pk"
	AppDir  = "pk"
	DBDir   = "db"
	CredDir = "creds"
)

var (
	ErrNotFound          = errors.New("not found")
	ErrCouldNotCreateAcc = errors.New("could not create account")
	ErrPermissionDenied  = errors.New("permission denied")
	ErrInternalError     = errors.New("internal error, possible db compromise")
	ErrCriticalFailure   = errors.New("could not perform critical operation")
)

type RequestDecoder interface {

}

type ResponseEncoder interface {
	Encode(response interface{})
}

// Hasher specifies an API for generating hashes of an arbitrary textual
// content.
type Hasher interface {
	// Hash generates the hashed string from plain-text.
	Hash(string) (string, error)

	// Compare compares plain-text version to the hashed one.
	//An error should indicate failed comparison.
	Compare(string, string) error
}

type Signer interface {
	Sign(string) ([]byte, []byte, error)

	Verify(password string, dbDigest []byte, dbSignature []byte) (err error)
}

type Encoder interface {
	Encode(password string) ([]byte, error)
	Decode(encoded []byte) (string, error)
}

type EncoderSigner interface {
	Encoder
	Signer
}

type BulkAddRequest struct {
	Token string `json:"token"`
	Accounts []Account `json:"accounts"`
}

type Account struct {
	Name     string `json:"name,omitempty"`
	UserName string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
	Created  string `json:"created,omitempty"`
}

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

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	Err   error  `json:"err"`
}

type AddRequest struct {
	Token    string `json:"token"`
	Name     string `json:"name"`
	UserName string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type GetRequest struct {
	Token    string `json:"token"`
	Name     string `json:"name"`
	UserName string `json:"username"`
}

type GetResponse struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Err error `json:"err"`
}

type ListResponse struct {
	Accounts []Account `json:"accounts"`
	Err      error     `json:"err,omitempty"`
}

type UpdateRequest struct {
	Token    string `json:"token"`
	Name     string `json:"name"`
	Username string `json:"username"`
	NewUser  string `json:"new_user"` //new username of the account
	Password string `json:"password"`
	Email    string `json:"email"`
}

//ErrResponse is a generic error response for function
//returning error as the only return value
type ErrResponse struct {
	Err error `json:"err"`
	Msg string `json:"msg"`
}

func (a Account) toDBAccount(keeper passwordKeeper) (DBAccount,error) {

	hash,err := keeper.hasher.Hash(a.Password)

	if err != nil {
		return DBAccount{}, err
	}

	encodedBytes,err := keeper.es.Encode(a.Password)

	if err != nil {
		return DBAccount{}, err
	}

	digestB, signB, err := keeper.es.Sign(a.Password)

	if err != nil {
		return DBAccount{}, err
	}

	return DBAccount{
		Name:      a.Name,
		UserName:  a.UserName,
		Email:     a.Email,
		Hash:      hash,
		Encoded:   encodedBytes,
		Digest:    digestB,
		Signature: signB,
		Created:   a.Created,
	},nil
}

func (a DBAccount) toAccount(keeper passwordKeeper) (Account,error) {

	pass,err:= keeper.es.Decode(a.Encoded)

	if err != nil {
		return Account{}, err
	}
	return Account{
		Name:     a.Name,
		UserName: a.UserName,
		Email:    a.Email,
		Password: pass,
		Created:  a.Created,
	},nil
}

type PasswordKeeper interface {

	//Register creates a new account. The function takes
	//Username, Email, Password. Account name is "master"
	//This is the first function to be called when running
	//the app
	Register(ctx context.Context, request RegisterRequest) (err ErrResponse)

	//Login function returns token after a user has supplied
	//his correct username and password or else an error
	//It also works fine if the email is supplied in place of
	//username
	Login(ctx context.Context, request LoginRequest) (response LoginResponse)

	//Add a new account. It takes token and new account details.
	//The method returns err if the process is not allowed
	Add(ctx context.Context, request AddRequest) (err ErrResponse)

	//Get returns details of the account after a user has supplied token
	//username and name of the account
	//username e.g pius
	//name e.g github
	Get(ctx context.Context, request GetRequest) (response GetResponse)

	//Delete removes the details of the account after a user has supplied token
	//username and name of the account
	//username e.g pius
	//name e.g github
	Delete(ctx context.Context, request GetRequest) (err ErrResponse)

	//List returns all the accounts registered under the master
	//accounts
	List(ctx context.Context) (list ListResponse)

	//Updates the details of the account
	Update(ctx context.Context, request UpdateRequest) (response ErrResponse)

	CredStore() CredStore

	//AddMany
	AddMany(ctx context.Context, req BulkAddRequest)(err error)
}

type PasswordStore interface {
	CheckAccount(ctx context.Context, name, username string) (err error)
	AddOwner(ctx context.Context, account Account) (err error)
	Add(ctx context.Context, account DBAccount) (err error)
	Get(ctx context.Context, name, username string) (account DBAccount, err error)
	GetOwner(ctx context.Context, name, username string) (account Account, err error)
	Delete(ctx context.Context, name, username string) (err error)
	Update(ctx context.Context, name, username string, account DBAccount) (err error)
	List(ctx context.Context) (accounts []DBAccount, err error)
}

type passwordKeeper struct {
	hasher    Hasher
	passwords PasswordStore
	tokenizer Tokenizer
	es        EncoderSigner
	credentials CredStore
}



func (p passwordKeeper) CredStore() CredStore {
	return p.credentials
}

var _ PasswordKeeper = (*passwordKeeper)(nil)

func NewPasswordKeeper(
	hasher Hasher, store PasswordStore,
	tokenizer Tokenizer, es EncoderSigner, cs CredStore) PasswordKeeper {
	return &passwordKeeper{
		hasher:    hasher,
		passwords: store,
		tokenizer: tokenizer,
		es:        es,
		credentials: cs,
	}
}

func (p passwordKeeper) AddMany(ctx context.Context, req BulkAddRequest) (err error) {
	tokenStr := req.Token
	//fixme: check the id in token and compare it to master
	_, err1 := p.tokenizer.Parse(tokenStr)
	if err1 != nil {
		return errors.Wrap(ErrPermissionDenied,err1)
	}

	accounts := req.Accounts


	for index, acc := range accounts{
		fmt.Printf("adding account no: %v\n",index+1)
		var a Account
		var d DBAccount
		now := time.Now().Format(time.RFC3339)
		name := acc.Name
		username:= acc.UserName
		password := acc.Password
		email := acc.Email

		a = Account{
			Name:     name,
			UserName: username,
			Email:    email,
			Password: password,
			Created:  now,
		}

		d,err = a.toDBAccount(p)

		if err != nil {
			return err
		}

		err = p.passwords.Add(ctx, d)

		if err != nil {
			return err
		}


	}

	return nil
}


func (p passwordKeeper) Register(ctx context.Context, request RegisterRequest) (errResponse ErrResponse) {

	password, err := p.hasher.Hash(request.Password)
	if err != nil {
		msg := "could not hash the password"
		return ErrResponse{
			Err: err,
			Msg: msg,
		}
	}

	email := request.Email
	username := request.Username

	created := time.Now().UTC().Format(time.RFC3339)
	dbAccount := Account{
		Name:     "master",
		UserName: username,
		Email:    email,
		Password: password,
		Created:  created,
	}

	err = p.passwords.AddOwner(ctx, dbAccount)

	if err != nil {
		msg := "could not add new user"
		return ErrResponse{
			Err: err,
			Msg: msg,
		}
	}

	return ErrResponse{Err: nil}
}

func (p passwordKeeper) Login(ctx context.Context, request LoginRequest) (response LoginResponse) {
	username := request.UserName
	password := request.Password
	account, err := p.passwords.GetOwner(ctx, "master", username)
	if err != nil {
		return LoginResponse{
			Token: "",
			Err:   err,
		}
	}

	err = p.hasher.Compare(password,account.Password)
	if err != nil {

		return LoginResponse{
			Token: "",
			Err:   err,
		}
	}

	token := NewToken(account.UserName)

	tokenStr, err := p.tokenizer.Issue(token)
	if err != nil {

		return LoginResponse{
			Token: "",
			Err:   err,
		}
	}

	return LoginResponse{
		Token: tokenStr,
		Err:   nil,
	}
}

func (p passwordKeeper) Add(ctx context.Context, request AddRequest) (err ErrResponse) {
	tokenStr := request.Token

	//fixme: check the id in token and compare it to master
	_, err1 := p.tokenizer.Parse(tokenStr)
	if err1 != nil {
		msg := "invalid token"
		return ErrResponse{
			Err: err1,
			Msg: msg,
		}
	}

	name := request.Name
	username:= request.UserName
	password := request.Password
	email := request.Email

	account := Account{
		Name:     name,
		UserName: username,
		Email:    email,
		Password: password,
		Created:  time.Now().Format(time.RFC3339),
	}

	dbAccount,err2 := account.toDBAccount(p)

	if err2 != nil {
		msg := "could not encrypt account details"
		return ErrResponse{
			Err: err2,
			Msg: msg,
		}
	}

	err3 := p.passwords.Add(ctx,dbAccount)

	if err3 != nil {
		msg := "could not store the account details"
		return ErrResponse{
			Err: err3,
			Msg: msg,
		}
	}


	return ErrResponse{Err: nil}
}

func (p passwordKeeper) Get(ctx context.Context, request GetRequest) (response GetResponse) {

	tokenStr := request.Token
	_,err :=p.tokenizer.Parse(tokenStr)

	if err != nil {
		return GetResponse{
			Email:    "",
			Password: "",
			Err: err,
		}
	}
	username := request.UserName
	name := request.Name
	account,err := p.passwords.Get(ctx,name,username)
	if err != nil {
		return GetResponse{
			Email:    "",
			Password: "",
			Err:      err,
		}
	}

	acc, err := account.toAccount(p)

	if err != nil {
		return GetResponse{
			Email:    "",
			Password: "",
			Err:      err,
		}
	}

	return GetResponse{
		Email:   acc.Email,
		Password: acc.Password,
		Err:      nil,
	}
}

func (p passwordKeeper) Delete(ctx context.Context, request GetRequest) (err ErrResponse) {
	panic("implement me")
}

func (p passwordKeeper) List(ctx context.Context) (list ListResponse) {
	panic("implement me")
}

func (p passwordKeeper) Update(ctx context.Context, request UpdateRequest) (response ErrResponse) {
	panic("implement me")
}
