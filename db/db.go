package db

import (
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

func DEX_ExampleDB_Model() {
	db := pg.Connect(&pg.Options{
		User:     "postgres",
		Password: "password",
		Database: "go-indexer",
	})
	defer db.Close()
	err := CreateSchema(db)
	if err != nil {
		panic(err)
	}
}

// createSchema creates database schema for User and Story models.
func CreateSchema(db *pg.DB) error {
	models := []interface{}{
		(*Token)(nil),
		(*Dex)(nil),
		(*Pool)(nil),
	}

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			Temp:        false,
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
