package indexer

import (
	"context"
	"encoding/binary"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
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

func GetPair(eth *ethclient.Client, tokenA, tokenB, factoryAddress string, blockNumber *big.Int) (pair_address common.Address, err error) {
	err = nil
	zeros := make([]byte, 12)
	data := make([]byte, 4+12)
	bytesA := common.HexToAddress(tokenA).Bytes()
	bytesB := common.HexToAddress(tokenB).Bytes()
	binary.BigEndian.PutUint32(data[:4], 0xe6a43905)
	data = append(data, bytesA...)
	data = append(data, zeros...)
	data = append(data, bytesB...)
	dexAddress := common.HexToAddress(factoryAddress)
	msg := ethereum.CallMsg{
		To:   &dexAddress,
		Data: data,
	}
	pairAddressBytes, err := eth.CallContract(context.Background(), msg, blockNumber)
	if err != nil {
		panic(err)
	}
	pair_address = common.BytesToAddress(pairAddressBytes)
	return
}

func GetReserves(eth *ethclient.Client, poolAddr common.Address, decimalsA uint8, decimalsB uint8, blockNumber *big.Int) (xVirtual, yVirtual float64, err error) {
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, 0x0902f1ac)
	msg := ethereum.CallMsg{
		To:   &poolAddr,
		Data: data,
	}
	resp, err := eth.CallContract(context.Background(), msg, blockNumber)
	if err != nil {
		return .0, .0, errors.Wrap(err, "failed to get reserves")
	}
	xVirtual = BigIntToFloat(new(big.Int).SetBytes(resp[:32]), decimalsA)
	yVirtual = BigIntToFloat(new(big.Int).SetBytes(resp[32:64]), decimalsB)
	return
}

func GetTokenInfo(eth *ethclient.Client, tokenAddress common.Address) (decimals uint8, symbol string, err error) {
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
	if len(decimalsBytes) > 0 {
		return decimalsBytes[len(decimalsBytes)-1], string(words), nil
	} else {
		return 0, string(words), nil
	}
}
