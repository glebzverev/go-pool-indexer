package arb

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/glebzverev/go-pool-indexer/indexer"
)

func (arb *Arb) runListenEventsLoop(eth *ethclient.Client) error {
	var poolAddresses []common.Address
	for poolAddress := range arb.pools {
		poolAddresses = append(poolAddresses, poolAddress)
	}

	SubscriptionTopics := []common.Hash{Topics.SyncUni}
	query := ethereum.FilterQuery{
		Addresses: poolAddresses,
		Topics:    [][]common.Hash{SubscriptionTopics},
	}
	ch := make(chan types.Log, 20)
	sub, err := eth.SubscribeFilterLogs(context.Background(), query, ch)
	if err != nil {
		panic(err)
	}
	go func() {
		defer sub.Unsubscribe()
		for {
			select {
			case err := <-sub.Err():
				log.Fatal(err)
			case parsedEvent := <-ch:
				timeEvent := time.Now()
				fmt.Println(timeEvent)
				arb.proceeSyncEvent(parsedEvent)
			}
		}
	}()
	return nil
}

func parseSyncEvent(event types.Log) (*big.Int, *big.Int) {
	return new(big.Int).SetBytes(event.Data[:32]), new(big.Int).SetBytes(event.Data[32:64])
}

func (arb *Arb) proceeSyncEvent(event types.Log) {
	arb.syncMutex.Lock()
	r1, r2 := parseSyncEvent(event)
	pool := arb.pools[event.Address]
	r1l := pool.reserve1
	r2l := pool.reserve2
	r1f := indexer.BigIntToFloat(r1, pool.token0.Decimals)
	r2f := indexer.BigIntToFloat(r2, pool.token1.Decimals)
	token0 := pool.token0.Address
	token1 := pool.token1.Address
	arb.linksMerge[token0][token1].cammulativeReserve1 += -r1l + r1f
	arb.linksMerge[token0][token1].cammulativeReserve2 += -r2l + r2f
	fmt.Printf("New reserves, %f %f for %s %s: <%s> \n", r1f, r2f, pool.token0.Symbol, pool.token1.Symbol, event.Address)
	arb.syncMutex.Unlock()

}

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
