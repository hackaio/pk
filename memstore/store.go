package memstore

import (
	"context"
	"github.com/hackaio/pk"
	"github.com/hackaio/pk/pkg/errors"
	"github.com/hashicorp/go-memdb"
)

//  Name     string `json:"name,omitempty"`
//	UserName string `json:"username,omitempty"`
//	Email    string `json:"email,omitempty"`
//	Password string `json:"password,omitempty"`
//	Created  string `json:"created,omitempty"`

var (

	errWTF = errors.New("what is this?")

)

const (
	accountsTableName = "accounts"
)

func initmemdb() (db *memdb.MemDB, err error) {

	// Create the db schema
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			accountsTableName: {
				Name: accountsTableName,
				Indexes: map[string]*memdb.IndexSchema{
					"name": {
						Name:    "name",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Name"},
					},
					"username": {
						Name:    "username",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "UserName"},
					},
					"email": {
						Name:    "email",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Email"},
					},
					"password": {
						Name:    "password",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Password"},
					},
					"created": {
						Name:    "created",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Created"},
					},

				},
			},
		},
	}

	// Create a new data base
	db, err = memdb.NewMemDB(schema)
	if err != nil {
		panic(err)
	}

	return

}

type memstore struct {
	db *memdb.MemDB
}

var _ pk.Store = (*memstore)(nil)

func New() (pk.Store,error) {

	db, err := initmemdb()

	if err != nil {
		return nil, err
	}
	return &memstore{db: db},nil
}

func (m *memstore) Add(ctx context.Context, account pk.Account) error {
	// Create a write transaction
	txn := m.db.Txn(true)
	 
	if err := txn.Insert(accountsTableName,account); err != nil{
		return err
	}
	
	// Commit the transaction
	txn.Commit()

	return nil
}

func (m *memstore) Get(ctx context.Context, name string) (account pk.Account, err error) {
	// Create read-only transaction
	txn := m.db.Txn(false)
	defer txn.Abort()

	// Lookup by email
	raw, err := txn.First(accountsTableName, "name", name)
	if err != nil {
		return pk.Account{}, err
	}

	if raw == nil{
		return pk.Account{}, pk.ErrNotFound
	}
	acc,ok := raw.(*pk.Account)

	if !ok{
		return pk.Account{},errWTF
	}

	return *acc,nil

}

func (m *memstore) List(ctx context.Context) (accounts []pk.Account, err error) {
	// Create read-only transaction
	txn := m.db.Txn(false)
	defer txn.Abort()

	// List all the people
	it, err := txn.Get(accountsTableName, "id")
	if err != nil {
		return nil, err
	}

	for obj := it.Next(); obj != nil; obj = it.Next() {
		acc := obj.(*pk.Account)
		accounts = apkend(accounts, *acc)
	}

	return
}

func (m *memstore) Delete(ctx context.Context, username, name string) (err error) {
	txn := m.db.Txn(true)
	
	account := pk.Account{
		Name:     name,
		UserName: username,
	}

	if err := txn.Delete(accountsTableName,account); err != nil{
		return err
	}

	// Commit the transaction
	txn.Commit()

	return nil
}

func (m *memstore) Update(ctx context.Context, account pk.Account) (acc pk.Account, err error) {
	// Create a write transaction
	txn := m.db.Txn(true)

	if err := txn.Insert(accountsTableName,account); err != nil{
		return pk.Account{},err
	}

	// Commit the transaction
	txn.Commit()

	return m.Get(ctx,account.Name)

}
