package pp

import (
	"context"
	"github.com/hackaio/pp/pkg/errors"
	"time"
)

type AccType int

const (
	Master AccType = iota
	Normal
)

var (
	ErrNotFound = errors.New("not found")
	ErrCouldNotCreateAcc = errors.New("could not create account")
)

// Hasher specifies an API for generating hashes of an arbitrary textual
// content.
type Hasher interface {
	// Hash generates the hashed string from plain-text.
	Hash(string) (string, error)

	// Compare compares plain-text version to the hashed one.
	//An error should indicate failed comparison.
	Compare(string, string) error
}

type Account struct {
	Name     string `json:"name,omitempty"`
	UserName string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
	Created  string `json:"created,omitempty"`
}

// Service specify an API for pp commandline tool
type Service interface {
	//Init initializes new account that multiple passwords
	//will be registered under it
	Init(ctx context.Context, username,email,password string) (err error)

	//Add create new acc
	Add(ctx context.Context,account Account)(err error)

	//Get returns a password of a specified account use id/username/name
	//of account
	//If you want to retrieve github password
	//username can be piusalfred
	//account name is github
	//since th username can be observed on multiple accs it is advised to
	//use name e.g github
	Get(ctx context.Context, username string) (acc Account, err error)

	//List retrieve all accounts
	List(ctx context.Context) (accounts []Account, err error)

	//Delete removes a specified account
	Delete(ctx context.Context, name, username string) (err error)

	//Update passwords
	Update(ctx context.Context, account Account) (acc Account, err error)
}

type Store interface {
	Add(ctx context.Context, account Account) error
	Get(ctx context.Context, name string) (account Account, err error)
	List(ctx context.Context) (accounts []Account, err error)
	Delete(ctx context.Context, username, name string) (err error)
	Update(ctx context.Context, account Account) (acc Account, err error)
}

type pp struct {
	store  Store
	hasher Hasher
}



var _ Service = (*pp)(nil)

func NewInstance(store Store, hasher Hasher) Service {
	return &pp{
		store:  store,
		hasher: hasher,
	}
}

func (p *pp) Init(ctx context.Context, username, email, password string) (err error) {
	hash,err := p.hasher.Hash(password)

	if err != nil {
		return err
	}

	now := time.Now().Format(time.RFC3339)
	account := Account{
		Name:     "master",
		UserName: username,
		Email:    email,
		Password: hash,
		Created:  now,
	}

	//lookup for the account
	_, err = p.store.Get(ctx, account.Name)
	if err != nil {
		if err == ErrNotFound{
			return p.store.Add(ctx,account)
		}
		return err
	}

	return ErrCouldNotCreateAcc
}

func (p *pp) Add(ctx context.Context, account Account) (err error) {
	hash,err := p.hasher.Hash(account.Password)

	if err != nil {
		return err
	}

	now := time.Now().Format(time.RFC3339)
	acc := Account{
		Name:     account.Name,
		UserName: account.UserName,
		Email:    account.Email,
		Password: hash,
		Created:  now,
	}

	//lookup for the account
	_, err = p.store.Get(ctx, acc.Name)
	if err != nil {
		if err == ErrNotFound{
			return p.store.Add(ctx,acc)
		}
		return err
	}

	return ErrCouldNotCreateAcc
}


func (p *pp) Get(ctx context.Context, name string) (acc Account, err error) {
	acc,err = p.store.Get(ctx,name)
	return
}

func (p *pp) List(ctx context.Context) (accounts []Account, err error) {
	accounts,err = p.store.List(ctx)
	return
}

func (p *pp) Delete(ctx context.Context, name, username string) (err error) {
	err = p.store.Delete(ctx,username,name)
	return
}

func (p *pp) Update(ctx context.Context, account Account) (acc Account, err error) {
	acc,err = p.store.Update(ctx,account)
	return
}

