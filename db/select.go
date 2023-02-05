package db

import (
	"github.com/go-pg/pg/v10"
)

func SelectTokens(db *pg.DB) []Token {
	tokens := new([]Token)
	err := db.Model(tokens).Select()
	if err != nil {
		panic(err)
	}
	return *tokens
}

func SelectDexes(db *pg.DB) (dexes []Dex) {
	err := db.Model(&dexes).Select()
	if err != nil {
		panic(err)
	}
	return
}

func SelectPools(db *pg.DB) (pools []Pool) {
	err := db.Model(&pools).Select()
	if err != nil {
		panic(err)
	}
	return
}
