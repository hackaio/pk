package pg

import (
	"context"
	"database/sql"
	"github.com/hackaio/pk"
)

type pgStore struct {
	db *sql.DB
}

var _ pk.PasswordStore = (*pgStore)(nil)

func NewStore(db *sql.DB) pk.PasswordStore {
	return pgStore{db: db}
}

func (p pgStore) AddOwner(ctx context.Context, account pk.Account) (err error) {
	panic("implement me")
}

func (p pgStore) GetOwner(ctx context.Context, name, username string) (account pk.Account, err error) {
	panic("implement me")
}

func (p pgStore) CheckAccount(ctx context.Context, name, username string) (err error) {
	panic("implement me")
}

func (p pgStore) Add(ctx context.Context, account pk.DBAccount) (err error) {
	panic("implement me")
}

func (p pgStore) Get(ctx context.Context, name, username string) (account pk.DBAccount, err error) {
	panic("implement me")
}

func (p pgStore) Delete(ctx context.Context, name, username string) (err error) {
	panic("implement me")
}

func (p pgStore) Update(ctx context.Context, name, username string, account pk.DBAccount) (err error) {
	panic("implement me")
}

func (p pgStore) List(ctx context.Context) (accounts []pk.DBAccount, err error) {
	panic("implement me")
}



