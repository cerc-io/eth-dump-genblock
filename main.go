// This script reads a geth genesis config and outputs the corresponding genesis block as
// RPC-compatible JSON.
//
// E.g.: eth-dump-genblock genesis.json > genesis_block.json
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/trie"
)

func main() {
	if len(os.Args) != 2 {
		fatal(fmt.Sprintf("Usage: %s genesis.json", os.Args[0]))
	}
	configPath := os.Args[1]

	file, err := os.Open(configPath)
	if err != nil {
		fatal("Failed to read genesis config file:", err)
	}

	genesis := new(core.Genesis)
	if err := json.NewDecoder(file).Decode(genesis); err != nil {
		fatal("Failed to parse genesis config:", err)
	}

	db := rawdb.NewMemoryDatabase()
	block, err := genesis.Commit(db, trie.NewDatabase(db))
	if err != nil {
		fatal("Failed to create genesis block:", err)
	}

	rpcBlock := RPCMarshalBlock(block)
	blockJSON, err := json.MarshalIndent(rpcBlock, "", "  ")
	if err != nil {
		fatal("Failed to convert block to JSON:", err)
	}

	fmt.Println(string(blockJSON))
}

func fatal(args ...any) {
	fmt.Fprintln(os.Stderr, args...)
	os.Exit(1)
}
