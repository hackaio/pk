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

package pg

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/hackaio/pk"
	"github.com/hackaio/pk/pkg/errors"
	"github.com/hackaio/pk/sql/stmt"
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

	//type Account struct {
	//	Name     string `json:"name,omitempty"`
	//	UserName string `json:"username,omitempty"`
	//	Email    string `json:"email,omitempty"`
	//	Password string `json:"password,omitempty"`
	//	Created  string `json:"created,omitempty"`
	//}
	//
	//type DBAccount struct {
	//	Name      string `json:"name,omitempty"`
	//	UserName  string `json:"username,omitempty"`
	//	Email     string `json:"email,omitempty"`
	//	Hash      []byte `json:"hash,omitempty"`
	//	Encoded   []byte `json:"encoded,omitempty"`
	//	Digest    []byte `json:"digest,omitempty"`
	//	Signature []byte `json:"signature,omitempty"`
	//	Created   string `json:"created,omitempty"`
	//}

	createMasterDb := `
CREATE TABLE IF NOT EXISTS masters(
    name VARCHAR (200) NOT NULL,
    username VARCHAR (200) NOT NULL,
    email VARCHAR(200) NOT NULL,
    password VARCHAR(200) NOT NULL,
    created VARCHAR(100) NOT NULL,
    PRIMARY KEY (name,username)
)
`

	createAccountsDb := `
CREATE TABLE IF NOT EXISTS accounts(
    name VARCHAR (200) NOT NULL,
    username VARCHAR (200) NOT NULL,
    email VARCHAR(200) NOT NULL,
    hash VARCHAR (300) NOT NULL UNIQUE,
    encoded VARCHAR (300) NOT NULL UNIQUE,
    digest VARCHAR (300) NOT NULL UNIQUE,
    signature VARCHAR (300) NOT NULL UNIQUE,
    created VARCHAR(100) NOT NULL,
    PRIMARY KEY (name,username)
)
`

	_, err = db.Exec(createMasterDb)
	_, err = db.Exec(createAccountsDb)

	if err != nil {
		errMsg := errors.New("could not create tables")
		return nil, errors.Wrap(err, errMsg)
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
	_, err = p.db.Exec(stmt.ADD_OWNER,account.Name, account.UserName,
		account.Email, account.Password, account.Created)
	return err
}

func (p pgStore) GetOwner(ctx context.Context, name, username string) (account pk.Account, err error) {

	err = p.db.QueryRow(stmt.GET_OWNER, name,username).Scan(&account.Name, &account.UserName,
		&account.Email, &account.Password, &account.Created)

	return account, err
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
