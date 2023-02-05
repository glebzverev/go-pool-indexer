package indexer

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/glebzverev/go-pool-indexer/db"
	"github.com/go-pg/pg/v10"
)

var (
	zeroAddress = common.HexToAddress("0x0000000000000000000000000000000000000000")
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func JsonToDataBase(eth *ethclient.Client, dataBase *pg.DB) {
	// eth, err := ethclient.Dial(os.Getenv("KEY"))
	// if err != nil {
	// 	panic(err)
	// }
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

func PoolsInit(eth *ethclient.Client, dataBase *pg.DB) {
	tokens := db.SelectTokens(dataBase)
	dexes := db.SelectDexes(dataBase)
	blockNumber, err := eth.BlockNumber(context.Background())
	check(err)
	for _, dex := range dexes {
		for i := 0; i < len(tokens)-1; i++ {
			for j := i; j < len(tokens); j++ {
				pairAddress, err := GetPair(eth, tokens[i].Address, tokens[j].Address, dex.FactoryAddress, nil)
				check(err)
				if pairAddress != zeroAddress {
					x, y, err := GetReserves(eth, pairAddress, nil)
					check(err)
					fmt.Println(x, y)
					reserves := &db.Reserves{
						Network:     dex.Network,
						Address:     pairAddress.String(),
						Reserve0:    x.String(),
						Reserve1:    y.String(),
						Liquidity:   new(big.Int).Mul(x, y).String(),
						BlockNumber: blockNumber,
					}
					pool := &db.Pool{
						Network:           dex.Network,
						DexInfo:           &dex,
						Address:           pairAddress.String(),
						Token0Address:     &tokens[i],
						Token1Address:     &tokens[j],
						LastReserveUpdate: reserves,
					}
					reserves.SafetyInsert(dataBase)
					pool.SafetyInsert(dataBase)
				}
			}
		}
	}

}
