package main

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/glebzverev/go-pool-indexer/db"
	"github.com/go-pg/pg/v10"
)

var (
	dataBase *pg.DB
)

func init() {
}

var (
	zeroAddress = common.HexToAddress("0x0000000000000000000000000000000000000000")
)

func main() {
	dataBase = pg.Connect(&pg.Options{
		User:     "postgres",
		Password: "password",
		Database: "go-indexer",
	})
	defer dataBase.Close()
	// db.CreateSchema(dataBase)

	// eth, err := ethclient.Dial(os.Getenv("ETH"))
	// if err != nil {
	// 	panic(err)
	// }

	// indexer.JsonToDataBase(eth, dataBase)
	// indexer.PoolsInit(eth, dataBase)

	fmt.Println(db.SelectResreves(dataBase))
	fmt.Println(db.SelectPools(dataBase))
}
