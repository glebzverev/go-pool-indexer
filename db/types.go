package db

import "fmt"

type Token struct {
	Network     string
	Address     string `pg:",pk"`
	Decimals    uint8
	Symbol      string
	TotalSupply float64
}

type TokenArray struct {
	Tokens  []Token
	current int
}

func (t Token) String() string {
	return fmt.Sprintf("\t%s <%d>\t", t.Symbol, t.Decimals)
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
	Reserve0    float64
	Reserve1    float64
	Liquidity   float64
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

func (r *TokenArray) Plan() int {
	return (len(r.Tokens))
}

func (r *TokenArray) Remaining() bool {
	// if r.current == r.Plan() {
	// 	r.current = 0
	// 	return false
	// }
	return r.current < len(r.Tokens)
}

func (r *TokenArray) Forward() (interface{}, error) {
	var token Token
	if r.current >= r.Plan() {
		return nil, fmt.Errorf("Last elem. Seg error")
	} else {
		token = r.Tokens[r.current]
		r.current++
	}
	return token, nil
}

type TokenInterface interface {
	String() string
}
