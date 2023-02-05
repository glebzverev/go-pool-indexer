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
