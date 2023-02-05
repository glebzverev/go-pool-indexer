package main

import (
	"fmt"

	"github.com/glebzverev/go-pool-indexer/indexer"
	"github.com/go-pg/pg/v10"
)

var (
	dataBase *pg.DB
)

func init() {
}

func main() {
	dataBase = pg.Connect(&pg.Options{
		User:     "postgres",
		Password: "password",
		Database: "go-indexer",
	})
	defer dataBase.Close()
	fmt.Println(indexer.Echo("Hello indexer"))
	indexer.JsonToDataBase(dataBase)
}