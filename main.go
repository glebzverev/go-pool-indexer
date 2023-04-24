package main

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/glebzverev/go-pool-indexer/arb"
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

// func main() {
// 	dataBase := pg.Connect(&pg.Options{
// 		User:     "diplomant",
// 		Password: "diplomant",
// 		Database: "diplom",
// 	})
// 	defer dataBase.Close()

// 	eth, err := ethclient.Dial(os.Getenv("ETH"))
// 	if err != nil {
// 		panic(err)
// 	}
// 	pools := db.SelectPools(dataBase)
// 	blockNumber, err := eth.BlockNumber(context.Background())
// 	if err != nil {
// 		panic(err)
// 	}
// 	blockNumberBig := new(big.Int).SetUint64(blockNumber)
// 	for _, pool := range pools {
// 		addr := common.HexToAddress(pool.Address)
// 		if pool.Token0Address == nil || pool.Token1Address == nil {
// 			continue
// 		}
// 		dec1 := pool.Token0Address.Decimals
// 		dec2 := pool.Token1Address.Decimals
// 		x, y, err := indexer.GetReserves(eth, addr, dec1, dec2, blockNumberBig)
// 		if err != nil {
// 			fmt.Println(err)
// 			continue
// 		}
// 		liq := x * y
// 		reserve := &db.Reserves{
// 			Network:     "ethereum",
// 			Address:     pool.Address,
// 			Reserve0:    x,
// 			Reserve1:    y,
// 			Liquidity:   liq,
// 			BlockNumber: blockNumber,
// 		}
// 		err = reserve.SafetyInsert(dataBase)
// 		if err != nil {
// 			panic(err)
// 		}
// 		// break
// 	}
// }

func main() {
	yamlPools := arb.ReadPairs()
	ARB := arb.CreateKnownPools(yamlPools)
	middles := ARB.GetMiddles(arb.TokenAddresses["USDT"], arb.TokenAddresses["WETH"])
	pool := ARB.GetPoolsByTokens(arb.TokenAddresses["USDT"], arb.TokenAddresses["WETH"])
	amountIn := 100.2
	token, res := pool.ComputeSwap(arb.TokenAddresses["USDT"], amountIn)
	fmt.Println(arb.Tokens[token].Symbol, res, amountIn/res)
	for _, middle := range middles {
		fmt.Println(arb.Tokens[middle].Symbol)
	}
	for {
	}
}

func parseSyncEvent(event *types.Log) (*big.Int, *big.Int) {
	return new(big.Int).SetBytes(event.Data[:32]), new(big.Int).SetBytes(event.Data[32:64])
}
