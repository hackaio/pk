package memstore

import (
	"context"
	"fmt"
	"github.com/hackaio/pp"
	"github.com/hashicorp/go-memdb"
)

func initmemdb()  (db *memdb.MemDB,err error){

	// Create the DB schema
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"person": &memdb.TableSchema{
				Name: "person",
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Email"},
					},
					"age": &memdb.IndexSchema{
						Name:    "age",
						Unique:  false,
						Indexer: &memdb.IntFieldIndex{Field: "Age"},
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
	DB *memdb.MemDB
}


var _ pp.Store = (*memstore)(nil)

func New() pp.Store {

	db,err := initmemdb()

	if err != nil {
		fmt.Printf("%v/n",err)
	}
	return &memstore{DB: db}
}

func (m *memstore) Add(ctx context.Context, account pp.Account) error {
	panic("implement me")
}

func (m *memstore) Get(ctx context.Context, name string) (account pp.Account, err error) {
	panic("implement me")
}

func (m *memstore) List(ctx context.Context) (accounts []pp.Account, err error) {
	panic("implement me")
}

func (m *memstore) Delete(ctx context.Context, username, name string) (err error) {
	panic("implement me")
}

func (m *memstore) Update(ctx context.Context, account pp.Account) (acc pp.Account, err error) {
	panic("implement me")
}
