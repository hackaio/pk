package pk

import (
	"context"
	"time"
)

// Service specify an API for pk commandline tool
type Service interface {

	//Init is the initialization process when pk is first run to set up
	//ownership and login details
	Init(ctx context.Context, username, email, password string) error

	//Auth take password and username and compares the password hash with the
	//one that was made during Init
	Auth(ctx context.Context, username, password string) (err error)

	//Add creates a new entry in the db if its not present yet
	Add(ctx context.Context, account Account) (err error)

	//Get retrieve the details of the password of account specified
	Get(ctx context.Context, name, username string) (account Account, err error)

	//List retrieve all of the details
	List(ctx context.Context) (accounts []Account, err error)

	//Delete
	Delete(ctx context.Context, name, username string) (err error)

	//Update
	Update(ctx context.Context, account Account) (acc Account, err error)
}

type service struct {
	store  Store
	hasher Hasher
}

var _ Service = (*service)(nil)

func NewInstance(store Store, hasher Hasher) Service {
	return &service{
		store:  store,
		hasher: hasher,
	}
}

// New returns a Service with all of the expected middleware wired in.
func NewPKService(store Store, hasher Hasher, middleware []Middleware) Service {
	var svc = NewInstance(store, hasher)
	for _, m := range middleware {
		svc = m(svc)
	}
	return svc
}

func (s service) Init(ctx context.Context, username, email, password string) error {
	hash, err := s.hasher.Hash(password)

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
	_, err = s.store.Get(ctx, account.Name, username)
	if err != nil {
		if err == ErrNotFound {
			return s.store.Add(ctx, account)
		}
		return err
	}

	return ErrCouldNotCreateAcc
}

func (s service) Auth(ctx context.Context, username, password string) (err error) {
	account, err := s.Get(ctx, "master", username)
	if err != nil {
		return err
	}

	err = s.hasher.Compare(password, account.Password)

	return
}

func (s service) Add(ctx context.Context, account Account) (err error) {
	hash, err := s.hasher.Hash(account.Password)

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
	_, err = s.store.Get(ctx, acc.Name, acc.UserName)
	if err != nil {
		if err == ErrNotFound {
			return s.store.Add(ctx, acc)
		}
		return err
	}

	return ErrCouldNotCreateAcc
}

func (s service) Get(ctx context.Context, name, username string) (account Account, err error) {
	account, err = s.store.Get(ctx, name, username)
	return
}

func (s service) List(ctx context.Context) (accounts []Account, err error) {
	accounts, err = s.store.List(ctx)
	return
}

func (s service) Delete(ctx context.Context, name, username string) (err error) {
	err = s.store.Delete(ctx, name, username)
	return
}

func (s service) Update(ctx context.Context, account Account) (acc Account, err error) {
	acc, err = s.store.Update(ctx, account)
	return
}
