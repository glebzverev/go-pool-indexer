package indexer

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

// var Selectors = struct {
// 	GetPair []byte
// }{
// 	GetPair: crypto.Keccak256([]byte("getPair(address,address)")),
// }
