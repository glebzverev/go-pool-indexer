package arb

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/glebzverev/go-pool-indexer/db"
	"github.com/glebzverev/go-pool-indexer/indexer"
	"github.com/go-pg/pg/v10"
	"gopkg.in/yaml.v3"
)

var TokenAddresses map[string]common.Address = map[string]common.Address{
	"USDT": common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7"),
	"WETH": common.HexToAddress("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"),
	"USDC": common.HexToAddress("0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48"),
	"BUSD": common.HexToAddress("0x4fabb145d64652a948d72533023f6e7a623c7c53"),
	"WBTC": common.HexToAddress("0x2260fac5e5542a773aa44fbcfedf7c193bc2c599"),
	"BNB":  common.HexToAddress("0xB8c77482e45F1F44dE1745F52C74426C631bDD52"),
}

var Tokens map[common.Address]*Token = map[common.Address]*Token{
	TokenAddresses["USDT"]: {
		Decimals: 6,
		Symbol:   "USDT",
		Address:  TokenAddresses["USDT"],
	},
	TokenAddresses["WETH"]: {
		Decimals: 18,
		Symbol:   "WETH",
		Address:  TokenAddresses["WETH"],
	},
	TokenAddresses["BUSD"]: {
		Decimals: 18,
		Symbol:   "BUSD",
		Address:  TokenAddresses["BUSD"],
	},
	TokenAddresses["USDC"]: {
		Decimals: 6,
		Symbol:   "USDC",
		Address:  TokenAddresses["USDC"],
	},
	TokenAddresses["WBTC"]: {
		Decimals: 8,
		Symbol:   "WBTC",
		Address:  TokenAddresses["WBTC"],
	},
	TokenAddresses["BNB"]: {
		Decimals: 18,
		Symbol:   "BNB",
		Address:  TokenAddresses["BNB"],
	},
}

type YamlPool struct {
	Token0Symbol string  `yaml:"token0"`
	Token1Symbol string  `yaml:"token1"`
	Fee          float64 `yaml:"fee"`
	Address      string  `yaml:"address"`
}
type YamlPools struct {
	Pools []YamlPool
}

func CheckPairs() {
	dataBase := pg.Connect(&pg.Options{
		User:     "diplomant",
		Password: "diplomant",
		Database: "diplom",
	})
	defer dataBase.Close()
	pools := db.SelectPools(dataBase)
	yamlPools := new(YamlPools)
	for _, pool := range pools {
		if pool.Token0Address != nil && pool.Token1Address != nil {
			token0Addr := common.HexToAddress(pool.Token0Address.Address)
			token1Addr := common.HexToAddress(pool.Token1Address.Address)
			_, ok0 := Tokens[token0Addr]
			_, ok1 := Tokens[token1Addr]
			if ok0 && ok1 {
				yamlPool := YamlPool{
					Token0Symbol: pool.Token0Address.Symbol,
					Token1Symbol: pool.Token1Address.Symbol,
					Address:      pool.Address,
					Fee:          0.003,
				}
				yamlPools.Pools = append(yamlPools.Pools, yamlPool)
			}
		}
	}
	yamlData, err := yaml.Marshal(yamlPools)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(yamlData))
}

func ReadPairs() YamlPools {
	f, err := os.ReadFile("data/pools.yml")
	if err != nil {
		log.Fatal(err)
	}
	var pools YamlPools
	if err := yaml.Unmarshal(f, &pools); err != nil {
		log.Fatal(err)
	}
	return pools
}

func CreateKnownPools(yampPools YamlPools) *Arb {
	arb := &Arb{
		syncMutex:  new(sync.Mutex),
		tokens:     Tokens,
		pools:      make(map[common.Address]*Pool, len(yampPools.Pools)),
		linksMerge: make(map[common.Address]map[common.Address]*ParallelPools, len(Tokens)),
		chains:     make(map[common.Address]map[common.Address][]common.Address, len(Tokens)),
	}
	eth, err := ethclient.Dial(os.Getenv("ETH_KEY"))
	if err != nil {
		panic(err)
	}
	for _, pool := range yampPools.Pools {
		tokenA := Tokens[TokenAddresses[pool.Token0Symbol]]
		tokenB := Tokens[TokenAddresses[pool.Token1Symbol]]
		if TokenAddresses[pool.Token0Symbol].Hex() > TokenAddresses[pool.Token1Symbol].Hex() {
			tokenB = Tokens[TokenAddresses[pool.Token0Symbol]]
			tokenA = Tokens[TokenAddresses[pool.Token1Symbol]]
		}
		knownPool := Pool{
			token0:  tokenA,
			token1:  tokenB,
			address: common.HexToAddress(pool.Address),
			fee:     pool.Fee,
		}
		xVirt, yVirt, err := indexer.GetReserves(
			eth,
			common.HexToAddress(pool.Address),
			uint8(knownPool.token0.Decimals),
			uint8(knownPool.token1.Decimals),
			nil,
		)
		if err != nil {
			panic(err)
		}
		knownPool.reserve1 = xVirt
		knownPool.reserve2 = yVirt
		arb.AddPool(&knownPool)
	}
	arb.ComputeChains()

	err = arb.runListenEventsLoop(eth)
	if err != nil {
		panic(err)
	}
	return arb
}
