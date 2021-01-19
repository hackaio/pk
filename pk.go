package pk

import (
	"context"
	"github.com/hackaio/pk/pkg/errors"
	"time"
)

const (
	AppDir   = "pk"
	DBDir    = "db"
	CredsDir = "creds"
)

var (
	ErrNotFound          = errors.New("not found")
	ErrCouldNotCreateAcc = errors.New("could not create account")
	ErrPermissionDenied  = errors.New("permission denied")
	ErrInternalError     = errors.New("internal error, possible db compromise")
	ErrCriticalFailure   = errors.New("could not perform critical operation")
)

type Store interface {
	Add(ctx context.Context, account Account) error
	Get(ctx context.Context, name, username string) (account Account, err error)
	List(ctx context.Context) (accounts []Account, err error)
	Delete(ctx context.Context, name, username string) (err error)
	Update(ctx context.Context, account Account) (acc Account, err error)
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

type Account struct {
	Name     string `json:"name,omitempty"`
	UserName string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
	Created  string `json:"created,omitempty"`
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
}

type ListResponse struct {
	Accounts []Account `json:"accounts"`
	Err      error     `json:"err"`
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
	Err string `json:"err"`
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
}

type PasswordStore interface {
	CheckAccount(ctx context.Context, name, username string) (err error)
	Add(ctx context.Context, account Account) (err error)
	Get(ctx context.Context, name, username string) (account Account, err error)
	Delete(ctx context.Context, name, username string) (err error)
	Update(ctx context.Context, name, username string, account Account) (err error)
	List(ctx context.Context) (accounts []Account, err error)
}

type passwordKeeper struct {
	hasher    Hasher
	passwords PasswordStore
	tokenizer Tokenizer
}

var _ PasswordKeeper = (*passwordKeeper)(nil)

func (p passwordKeeper) Register(ctx context.Context, request RegisterRequest) (errResponse ErrResponse) {
	password, err := p.hasher.Hash(request.Password)
	if err != nil {
		return ErrResponse{Err: err.Error()}
	}

	created := time.Now().UTC().Format(time.RFC3339)
	dbAccount := Account{
		Name:     "master",
		UserName: request.Username,
		Email:    request.Email,
		Password: password,
		Created:  created,
	}

	err = p.passwords.Add(ctx, dbAccount)

	if err != nil {
		return ErrResponse{Err: err.Error()}
	}

	return ErrResponse{}
}

func (p passwordKeeper) Login(ctx context.Context, request LoginRequest) (response LoginResponse) {
	username := request.UserName
	password := request.Password
	account, err := p.passwords.Get(ctx, "master", username)
	if err != nil {
		return LoginResponse{
			Token: "",
			Err:   err,
		}
	}

	err = p.hasher.Compare(account.Password, password)
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
	//token := request.Token
	return ErrResponse{}
}

func (p passwordKeeper) Get(ctx context.Context, request GetRequest) (response GetResponse) {
	panic("implement me")
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