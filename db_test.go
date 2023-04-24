package main

import (
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/glebzverev/go-pool-indexer/arb"
	"github.com/glebzverev/go-pool-indexer/db"
	"github.com/go-pg/pg/v10"
	"github.com/sbwhitecap/tqdm"
)

func TestSchema(t *testing.T) {
	dataBase := pg.Connect(&pg.Options{
		User:     "diplomant",
		Password: "diplomant",
		Database: "diplom",
	})
	defer dataBase.Close()
	fmt.Println(len(db.SelectTokens(dataBase)))
	fmt.Println(len(db.SelectPools(dataBase)))
	fmt.Println(len(db.SelectDexes(dataBase)))
}

func TestTqdm(t *testing.T) {
	dataBase := pg.Connect(&pg.Options{
		User:     "diplomant",
		Password: "diplomant",
		Database: "diplom",
	})
	defer dataBase.Close()
	fmt.Println(len(db.SelectTokens(dataBase)))
	fmt.Println(len(db.SelectPools(dataBase)))
	fmt.Println(len(db.SelectDexes(dataBase)))
	tokenArray := &db.TokenArray{
		Tokens: db.SelectTokens(dataBase),
	}

	err := tqdm.With(tokenArray, "hello", processToken[db.Token])
	if err != nil {
		panic(err)
	}
}

func TestArbAddress(t *testing.T) {
	arb.ReadPairs()
}

func processToken[T Type](v interface{}) (brk bool) {
	time.Sleep(time.Millisecond * 10)
	elem := v.(T)
	_, err := io.WriteString(os.Stdout, elem.String())
	if err != nil {
		panic(err)
	}
	return
}

type Type interface {
	db.Token | db.Pool | db.Dex

	String() string
}
