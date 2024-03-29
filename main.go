package main

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/glebzverev/go-pool-indexer/arb"
	"github.com/glebzverev/go-pool-indexer/server"
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

func main() {
	yamlPools := arb.ReadPairs()
	ARB := arb.CreateKnownPools(yamlPools)
	// middles := ARB.GetMiddles(arb.TokenAddresses["USDT"], arb.TokenAddresses["WETH"])
	pool := ARB.GetPoolsByTokens(arb.TokenAddresses["USDT"], arb.TokenAddresses["WETH"])
	amountIn := 100.2
	token, res := pool.ComputeSwap(arb.TokenAddresses["USDT"], amountIn)
	fmt.Println(arb.Tokens[token].Symbol, res, amountIn/res)
	opt, _ := ARB.FindOptimal(arb.TokenAddresses["USDT"], arb.TokenAddresses["WETH"], 1e+6)
	fmt.Println(opt)
	server.New(ARB)
	for {
	}
}
