package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"sort"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/glebzverev/go-pool-indexer/db"
	"github.com/glebzverev/go-pool-indexer/indexer"
	"github.com/go-pg/pg/v10"
	"github.com/sbwhitecap/tqdm"
)

var Topics = struct {
	Swap     common.Hash
	Mint     common.Hash
	Burn     common.Hash
	SyncUni  common.Hash
	SyncVelo common.Hash
	Transfer common.Hash
}{
	Swap:     crypto.Keccak256Hash([]byte("Swap(address,address,int256,int256,uint160,uint128,int24)")),
	Mint:     crypto.Keccak256Hash([]byte("Mint(address,address,int24,int24,uint128,uint256,uint256)")),
	Burn:     crypto.Keccak256Hash([]byte("Burn(address,int24,int24,uint128,uint256,uint256)")),
	SyncUni:  crypto.Keccak256Hash([]byte("Sync(uint112,uint112)")),
	SyncVelo: crypto.Keccak256Hash([]byte("Sync(uint256,uint256)")),
	Transfer: crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)")),
}
var (
	dataBase *pg.DB
)

func init() {
}

var (
	zeroAddress = common.HexToAddress("0x0000000000000000000000000000000000000000")
)

func existsInSlice[T comparable](val T, values []T) bool {
	for _, v := range values {
		if val == v {
			return true
		}
	}
	return false
}

func main3() {
	eth, err := ethclient.Dial(os.Getenv("ETH"))
	if err != nil {
		panic(err)
	}

	database := pg.Connect(&pg.Options{
		User:     "diplomant",
		Password: "diplomant",
		Database: "diplom",
	})
	defer database.Close()

	if createSchema {
		db.CreateSchema(database)
		tqdm.R(0, 10, func(v interface{}) (brk bool) {
			time.Sleep(1000 * time.Millisecond)
			return
		})
	}
	if tokenListen {
		blockNumber, err := eth.BlockNumber(context.Background())
		if err != nil {
			panic(err)
		}

		query := ethereum.FilterQuery{
			FromBlock: new(big.Int).SetUint64(blockNumber - 100),
			ToBlock:   new(big.Int).SetUint64(blockNumber - 1),
			Topics:    [][]common.Hash{{Topics.Transfer}},
		}

		logs, err := eth.FilterLogs(context.Background(), query)
		tokenAdresses := make(map[common.Address]int)
		if err != nil {
			log.Fatal(err)
		}
		if len(logs) == 0 {
			fmt.Println("Have no events to this story period")
		} else {
			for _, log := range logs {
				tokenAdresses[log.Address] += 1
			}
		}
		keys := make([]common.Address, 0, len(tokenAdresses))
		for key := range tokenAdresses {
			keys = append(keys, key)
		}

		sort.SliceStable(keys, func(i, j int) bool {
			return tokenAdresses[keys[i]] > tokenAdresses[keys[j]]
		})
		for _, k := range keys {
			if tokenAdresses[k] > 5 {

				decimals, symbol, err := indexer.GetTokenInfo(eth, k)
				if err != nil {
					fmt.Println(err)
				} else {
					token := &db.Token{
						Network:  "ethereum",
						Address:  k.String(),
						Decimals: decimals,
						Symbol:   symbol,
					}
					token.SafetyInsert(database)
				}
			}
		}
	}
	if jsonToDatabase {
		indexer.JsonToDataBase(eth, database)
		tqdm.R(0, 10, func(v interface{}) (brk bool) {
			time.Sleep(1000 * time.Millisecond)
			return
		})
	}
	if poolIndex {
		indexer.PoolsInit(eth, database)
		tqdm.R(0, 10, func(v interface{}) (brk bool) {
			time.Sleep(1000 * time.Millisecond)
			return
		})
	}
}

func main() {
	eth, err := ethclient.Dial(os.Getenv("ETH"))
	if err != nil {
		panic(err)
	}

	database := pg.Connect(&pg.Options{
		User:     "diplomant",
		Password: "diplomant",
		Database: "diplom",
	})
	defer database.Close()

	indexer.SyncListen(eth, database)
}

var (
	jsonToDatabase bool = true
	createSchema   bool = true
	poolIndex      bool = true
	tokenListen    bool = true
)
