package pg

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/hackaio/pk"
	_ "github.com/lib/pq"
)

const (
	hostname = "localhost"
	port     = "5432"
	user     = "postgres"
	password = "postgres"
	dbname   = "postgres"
	sslmode  = "disable"
)

// Connect creates a connection to the PostgreSQL instance and applies any
func Connect() (*sql.DB, error) {
	url := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", hostname, port, user, dbname, password, sslmode)

	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	return db, nil
}

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
