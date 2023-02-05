package db

import "fmt"

type Token struct {
	Network     string
	Address     string `pg:",pk"`
	Decimals    uint8
	Symbol      string
	TotalSupply string
}

func (t Token) String() string {
	return fmt.Sprintf("%s<%s %s %d>", t.Symbol, t.Network, t.Address, t.Decimals)
}

type Dex struct {
	Network        string
	Name           string
	FactoryAddress string `pg:",pk"`
	RouterAddress  string
}

func (d Dex) String() string {
	return fmt.Sprintf("%s<%s %s %s>", d.Name, d.Network, d.FactoryAddress, d.RouterAddress)
}

type Reserves struct {
	Network     string
	Address     string
	Reserve0    string
	Reserve1    string
	Liquidity   string
	BlockNumber uint64 // Reserve0 * Reserve1
}

type Pool struct {
	Network           string
	DexInfo           *Dex
	Address           string `pg:",pk"`
	Token0Address     *Token
	Token1Address     *Token
	LastReserveUpdate *Reserves
}

func (p Pool) String() string {
	return fmt.Sprintf("%s: %s <%s>", p.Network, p.DexInfo, p.Address)
}
