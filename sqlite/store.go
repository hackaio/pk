package sqlite

import (
	"context"
	"database/sql"
	"github.com/hackaio/pk"
	"github.com/hackaio/pk/sql/stmt"
)

type store struct {
	db *sql.DB
}

func NewStore(dbpath string) (s pk.Store, err error) {
	db, err := initdb(dbpath)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}
	return store{db: db}, nil
}

func (s store) Add(ctx context.Context, account pk.Account) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	st, err := tx.Prepare(stmt.ADD)
	if err != nil {
		return err
	}
	defer st.Close()

	name := account.Name
	uname := account.UserName
	email := account.Email
	password := account.Password
	created := account.Created

	_, err = st.Exec(name, uname, email, password, created)

	if err != nil {
		return err
	}
	err = tx.Commit()

	if err != nil {
		return err
	}

	return nil
}

func (s store) Get(ctx context.Context, name, username string) (account pk.Account, err error) {
	st, err := s.db.Prepare(stmt.GET)
	if err != nil {
		return pk.Account{}, err
	}
	defer st.Close()

	err = st.QueryRow(name, username).
		Scan(&account.Name, &account.UserName, &account.Email, &account.Password, &account.Created)
	if err != nil {
		return pk.Account{}, err
	}
	return
}

func (s store) List(ctx context.Context) (accounts []pk.Account, err error) {
	rows, err := s.db.Query(stmt.LIST)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		account := pk.Account{}
		err := rows.Scan(&account.Name, &account.UserName, &account.Email, &account.Password, &account.Created)
		if err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return accounts, nil

}

func (s store) Delete(ctx context.Context, name, username string) (err error) {
	// delete
	st, err := s.db.Prepare(stmt.DELETE)
	if err != nil {
		return err
	}

	_, err = st.Exec(name,username)
	if err != nil {
		return err
	}

	err = s.db.Close()

	return
}

func (s store) Update(ctx context.Context, account pk.Account) (acc pk.Account, err error) {
	panic("implement me")
}
