package arb

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Token struct {
	Decimals uint64
	Symbol   string
	Address  common.Address
}

type Pool struct {
	token0   *Token
	token1   *Token
	fee      *big.Int
	reserve1 *big.Int
	reserve2 *big.Int
}

type ParallelPools struct {
	pools               map[common.Address]*Pool
	token0              *Token
	token1              *Token
	cammulativeReserve0 *big.Int
	cammulativeReserve1 *big.Int
	cammulativeFee      *big.Int
}

type Arb struct {
	linksMerge map[common.Address]map[common.Address]*ParallelPools
	tokens     map[common.Address]*Token
	pools      map[common.Address]*Pool
}

func (arb *Arb) GetPoolsByTokens(token0 common.Address, token1 common.Address) ParallelPools {
	return *arb.linksMerge[token0][token1]
}

func (arb *Arb) GetPool(poolAddress common.Address) Pool {
	return *arb.pools[poolAddress]
}

func (arb *Arb) GetTokensInfo(tokenAddress common.Address) Token {
	return *arb.tokens[tokenAddress]
}

func (pp *ParallelPools) ComputeSwap(address common.Address, amountIn *big.Int) (common.Address, *big.Int) {
	if address == pp.token0.Address {
		val := new(big.Int)
		return pp.token1.Address, val
	} else if address == pp.token1.Address {
		val := new(big.Int)
		return pp.token0.Address, val
	} else {
		panic("Incorrect token for pool")
	}
}
