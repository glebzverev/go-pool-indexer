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
				fmt.Println(parseSyncEvent(parsedEvent))
				fmt.Println(timeEvent, parsedEvent)
			}
		}
	}()
	return nil
}

func parseSyncEvent(event types.Log) (*big.Int, *big.Int) {
	return new(big.Int).SetBytes(event.Data[:32]), new(big.Int).SetBytes(event.Data[32:64])
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
