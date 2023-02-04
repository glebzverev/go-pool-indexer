package main

import (
	"fmt"

	"github.com/glebzverev/go-pool-indexer/indexer"
)

func main() {
	fmt.Println(indexer.Echo("Hello indexer"))
	indexer.Schema()
}
