package indexer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/glebzverev/go-pool-indexer/db"
	"github.com/go-pg/pg/v10"
)

func JsonToDataBase(dataBase *pg.DB) {
	eth, err := ethclient.Dial(os.Getenv("KEY"))
	if err != nil {
		panic(err)
	}
	jsonFile, err := os.Open("./config.json")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened users.json")
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var chains Chains

	json.Unmarshal(byteValue, &chains)
	for _, chain := range chains.Chains {
		fmt.Println("Chain Name: " + chain.Name)
		for _, token := range chain.Tokens {
			decimals, symbol, err := getTokenInfo(eth, common.HexToAddress(token))
			if err != nil {
				fmt.Println("can't get token info", err)
			}
			fmt.Println(decimals, symbol)
			dbToken := &db.Token{
				Network:     chain.Name,
				Address:     token,
				Decimals:    decimals,
				Symbol:      symbol,
				TotalSupply: "10",
			}
			err = dbToken.SafetyInsert(dataBase)
			if err != nil {
				panic(err)
			}
		}
		for _, dex := range chain.Dexes {
			dbDex := &db.Dex{
				Network:        chain.Name,
				FactoryAddress: dex.Factory,
				Name:           dex.Name,
			}
			err = dbDex.SafetyInsert(dataBase)
			if err != nil {
				panic(err)
			}
		}
	}
}
