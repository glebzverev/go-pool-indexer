package main

import (
	"fmt"
	"testing"

	"github.com/glebzverev/go-pool-indexer/db"
	"github.com/go-pg/pg/v10"
)

func TestSchema(t *testing.T) {
	dataBase = pg.Connect(&pg.Options{
		User:     "postgres",
		Password: "password",
		Database: "go-indexer",
	})
	defer dataBase.Close()
	fmt.Println(db.SelectTokens(dataBase))
	fmt.Println(db.SelectDexes(dataBase))
}
