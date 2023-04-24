package indexer

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/glebzverev/go-pool-indexer/db"
	"github.com/go-pg/pg/v10"
)

func SyncListen(client *ethclient.Client, database *pg.DB) {
	dexes := db.SelectPools(database)
	var poolAddresses []common.Address
	for _, dex := range dexes {
		poolAddresses = append(poolAddresses, common.HexToAddress(dex.Address))
	}

	SubscriptionTopics := [][]common.Hash{{
		common.Hash(crypto.Keccak256Hash([]byte("Sync(uint112,uint112)"))),
		common.Hash(crypto.Keccak256Hash([]byte("Sync(uint256,uint256)"))),
	}}
	query := ethereum.FilterQuery{
		Addresses: poolAddresses,
		Topics:    SubscriptionTopics,
	}
	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatal(err)
	}
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLog := <-logs:
			processSync(database, &vLog) // pointer to event log
		}
	}
}

func processSync(database *pg.DB, event *types.Log) {
	xVirtual, yVirtual := parseSyncEvent(event)
	pool := new(db.Pool)
	db.GetPool(database, event.Address.String(), pool)
	x := (BigIntToFloat(xVirtual, pool.Token0Address.Decimals))
	y := (BigIntToFloat(yVirtual, pool.Token1Address.Decimals))
	reserves := db.Reserves{
		Network:     "ethereum",
		Address:     event.Address.String(),
		Reserve0:    x,
		Reserve1:    y,
		Liquidity:   x * y,
		BlockNumber: event.BlockNumber,
	}
	pool.LastReserveUpdate = &reserves
	err := pool.SafetyUpdate(database)
	fmt.Println("ERROR ", err)
	fmt.Println("SYNC", xVirtual, yVirtual)
}

func parseSyncEvent(event *types.Log) (*big.Int, *big.Int) {
	return new(big.Int).SetBytes(event.Data[:32]), new(big.Int).SetBytes(event.Data[32:64])
}
