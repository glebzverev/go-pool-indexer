package indexer

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
)

func Echo(s string) string {
	return s
}

func Schema() {
	eth, err := ethclient.Dial(os.Getenv("KEY"))
	if err != nil {
		panic(err)
	}
	jsonFile, err := os.Open("./config.json")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened users.json")
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var chains Chains

	json.Unmarshal(byteValue, &chains)
	for _, chain := range chains.Chains {
		fmt.Println("Chain Name: " + chain.Name)

		for i := 0; i < len(chain.Tokens)-1; i++ {
			for j := i + 1; j < len(chain.Tokens); j++ {

				getPair(eth, chain.Tokens[i], chain.Tokens[j], chain.Dexes)
			}
		}
	}
}

func getPair(eth *ethclient.Client, tokenA, tokenB string, Dexes []Dex) (pair_address common.Address, error err) {
	dataDecimals := make([]byte, 4)

	binary.BigEndian.PutUint32(dataDecimals, 0x313ce567)

	msg := ethereum.CallMsg{
		To:   &tokenAddress,
		Data: dataDecimals,
	}
	decimalsBytes, err := eth.CallContract(context.Background(), msg, nil)
	if err != nil {
		return "", errors.Wrap(err, "failed to get token decimals")
	}
	pair_address = common.HexToAddress(tokenA)
	fmt.Println(decimalsBytes)
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

type Chain struct {
	Name   string   `json:"name"`
	Tokens []string `json:"token_addresses"`
	Dexes  []Dex    `json:"dexes"`
}

type Chains struct {
	Chains []Chain `json:"chains"`
}

type Dex struct {
	Name    string `json:"name"`
	Factory string `json:"address"`
}
