package indexer

import (
	"context"
	"encoding/binary"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/glebzverev/go-pool-indexer/db"
	"github.com/pkg/errors"
)

func Echo(s string) string {
	return s
}

func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	pending := big.NewInt(-1)
	if number.Cmp(pending) == 0 {
		return "pending"
	}
	return hexutil.EncodeBig(number)
}

func GetPair(eth *ethclient.Client, tokenA, tokenB string, dex db.Dex, blockNumber *big.Int) (pair_address common.Address, err error) {
	var batch []rpc.BatchElem
	err = nil

	arg := toBlockNumArg(blockNumber)

	data := make([]byte, 4+12)

	binary.BigEndian.PutUint32(data[:4], 0xf30dba93)
	addressA := common.HexToAddress(tokenA)
	addressB := common.HexToAddress(tokenB)
	data = append(data, addressA.Bytes()...)
	zeros := make([]byte, 12)
	data = append(data, zeros...)
	data = append(data, addressB.Bytes()...)
	dexAddress := common.HexToAddress(dex.FactoryAddress)
	msg := ethereum.CallMsg{
		To:   &tokenAddress,
		Data: dataDecimals,
	}

			"data": hexutil.Bytes(data),
	}

	// decimalsBytes, err := eth.CallContract(context.Background(), msg, nil)
	// if err != nil {
	// 	return "", errors.Wrap(err, "failed to get token decimals")
	// }
	pair_address = common.HexToAddress(tokenA)
	// fmt.Println(decimalsBytes)
	return
}

func getTokenInfo(eth *ethclient.Client, tokenAddress common.Address) (decimals uint8, symbol string, err error) {
	dataDecimals := make([]byte, 4)
	dataSymbol := make([]byte, 4)

	binary.BigEndian.PutUint32(dataDecimals, 0x313ce567)
	msg := ethereum.CallMsg{
		To:   &tokenAddress,
		Data: dataDecimals,
	}
	decimalsBytes, err := eth.CallContract(context.Background(), msg, nil)
	if err != nil {
		return 0, "", errors.Wrap(err, "failed to get token decimals")
	}
	binary.BigEndian.PutUint32(dataSymbol, 0x95d89b41)
	msg = ethereum.CallMsg{
		To:   &tokenAddress,
		Data: dataSymbol,
	}
	dataSymbol, err = eth.CallContract(context.Background(), msg, nil)
	if err != nil {
		return 0, "", errors.Wrap(err, "failed to get token symbol")
	}
	words := []byte{}
	for i := 64; i < len(dataSymbol)-1; i += 1 {
		if binary.LittleEndian.Uint16(dataSymbol[i:]) == 0 {
			continue
		}
		words = append(words, dataSymbol[i])
	}
	return decimalsBytes[len(decimalsBytes)-1], string(words), nil
}
