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
