package arb

import (
	"math"

	"github.com/ethereum/go-ethereum/common"
)

type Token struct {
	Decimals uint8
	Symbol   string
	Address  common.Address
}

type Pool struct {
	token0   *Token
	token1   *Token
	address  common.Address
	fee      float64
	reserve1 float64
	reserve2 float64
}

type ParallelPools struct {
	pools               map[common.Address]*Pool
	token0              *Token
	token1              *Token
	cammulativeReserve1 float64
	cammulativeReserve2 float64
	cammulativeFee      float64
}

type Arb struct {
	linksMerge map[common.Address]map[common.Address]*ParallelPools
	tokens     map[common.Address]*Token
	pools      map[common.Address]*Pool
	chains     map[common.Address]map[common.Address][]common.Address
}

type WeightedPool struct {
	Pool  *Pool
	Alpha float64
}

func (pp *ParallelPools) DividedRoute() []WeightedPool {
	wp := make([]WeightedPool, 0)
	for _, pool := range pp.pools {
		alpha := math.Sqrt(pool.reserve1 / pp.cammulativeReserve1 * pool.reserve2 / pp.cammulativeReserve2)
		wp = append(wp, WeightedPool{
			Pool:  pool,
			Alpha: alpha,
		})
	}
	return wp
}

func (p *Pool) toParallelPool() *ParallelPools {
	return &ParallelPools{
		pools: map[common.Address]*Pool{
			p.address: p,
		},
		token0:              p.token0,
		token1:              p.token1,
		cammulativeReserve1: p.reserve1,
		cammulativeReserve2: p.reserve2,
		cammulativeFee:      p.fee,
	}
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

func (pp *ParallelPools) ComputeSwap(address common.Address, amountIn float64) (common.Address, float64) {
	if address == pp.token0.Address {
		val := pp.cammulativeReserve2 - pp.cammulativeReserve1*pp.cammulativeReserve2/(pp.cammulativeReserve1+(1-pp.cammulativeFee)*amountIn)
		return pp.token1.Address, val
	} else if address == pp.token1.Address {
		val := pp.cammulativeReserve1 - pp.cammulativeReserve1*pp.cammulativeReserve2/(pp.cammulativeReserve2+(1-pp.cammulativeFee)*amountIn)
		return pp.token0.Address, val
	} else {
		panic("Incorrect token for pool")
	}
}

func (pp *ParallelPools) AddPool(pool *Pool) {
	alpha := (pp.cammulativeReserve1/pool.reserve1 + pp.cammulativeReserve2/pool.reserve2) / 2
	if pool.token0.Address == pp.token0.Address {
		pp.cammulativeReserve1 += pool.reserve1
		pp.cammulativeReserve2 += pool.reserve2
	} else {
		pp.cammulativeReserve1 += pool.reserve2
		pp.cammulativeReserve2 += pool.reserve1
	}
	pp.pools[pool.address] = pool
	pp.cammulativeFee = (pp.cammulativeFee + alpha*pool.fee) / (1 + alpha)
}

func (arb *Arb) AddPool(pool *Pool) {
	if arb.pools[pool.address] != nil {
		return
	}
	arb.pools[pool.address] = pool
	if arb.linksMerge[pool.token0.Address] != nil {
		if arb.linksMerge[pool.token0.Address][pool.token1.Address] != nil {
			arb.linksMerge[pool.token0.Address][pool.token1.Address].AddPool(pool)
		} else {
			arb.linksMerge[pool.token0.Address][pool.token1.Address] = pool.toParallelPool()
		}
	} else {
		arb.linksMerge[pool.token0.Address] = map[common.Address]*ParallelPools{
			pool.token1.Address: pool.toParallelPool(),
		}
	}
	if arb.linksMerge[pool.token1.Address] != nil {
		if arb.linksMerge[pool.token1.Address][pool.token0.Address] != nil {
			arb.linksMerge[pool.token1.Address][pool.token0.Address].AddPool(pool)
		} else {
			arb.linksMerge[pool.token1.Address][pool.token0.Address] = pool.toParallelPool()
		}
	} else {
		arb.linksMerge[pool.token1.Address] = map[common.Address]*ParallelPools{
			pool.token0.Address: pool.toParallelPool(),
		}
	}

}

func (arb *Arb) addChain(start common.Address, middle common.Address, finish common.Address) {
	if _, ok := arb.chains[start]; ok {
		if _, ok = arb.chains[start][finish]; ok {
			if !include(arb.chains[start][finish], middle) {
				arb.chains[start][finish] = append(arb.chains[start][finish], middle)
			}
		} else {
			arb.chains[start][finish] = []common.Address{middle}
		}
	} else {
		middleArr := make([]common.Address, 1)
		middleArr[0] = middle
		arb.chains[start] = map[common.Address][]common.Address{
			finish: middleArr,
		}
	}
	if _, ok := arb.chains[finish]; ok {
		if _, ok = arb.chains[finish][start]; ok {
			if !include(arb.chains[finish][start], middle) {
				arb.chains[finish][start] = append(arb.chains[finish][start], middle)
			}
		} else {
			arb.chains[finish][start] = []common.Address{middle}
		}
	} else {
		middleArr := make([]common.Address, 1)
		middleArr[0] = middle
		arb.chains[finish] = map[common.Address][]common.Address{
			start: middleArr,
		}
	}

}

func (arb *Arb) ComputeChains() {
	for start := range arb.tokens {
		for middle := range arb.linksMerge[start] {
			for finish := range arb.linksMerge[middle] {
				if finish != start {
					arb.addChain(start, middle, finish)
				}
			}
		}
	}
}

func (arb *Arb) GetMiddles(start, finish common.Address) []common.Address {
	return arb.chains[start][finish]
}

func include[T common.Address](array []T, elem T) bool {
	for _, a := range array {
		if elem == a {
			return true
		}
	}
	return false
}
