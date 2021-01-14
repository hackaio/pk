package pp

import (
	"context"
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

// PP specify an API for pp commandline tool
type PP interface {
	//Init initializes new account that multiple passwords
	//will be registered under it
	Init(ctx context.Context, account Account) (err error)

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

var _ PP = (*pp)(nil)

func NewPPInstance(store Store, hasher Hasher) PP {
	return &pp{
		store:  store,
		hasher: hasher,
	}
}

func (p *pp) Init(ctx context.Context, account Account) (err error) {
	panic("implement me")
}

func (p *pp) Get(ctx context.Context, username string) (acc Account, err error) {
	panic("implement me")
}

func (p *pp) List(ctx context.Context) (accounts []Account, err error) {
	panic("implement me")
}

func (p *pp) Delete(ctx context.Context, name, username string) (err error) {
	panic("implement me")
}

func (p *pp) Update(ctx context.Context, account Account) (acc Account, err error) {
	panic("implement me")
}
