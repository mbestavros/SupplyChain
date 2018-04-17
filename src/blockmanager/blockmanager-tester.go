package blockmanager

import (
	"fmt"
	"time"
)

func main() {
	// initialize block manager
	bm := Blockmanager{}

	fmt.Println("Generating a genesis block (CREATE transaction)...")

	t := time.Now()
	genesisBlock := Block{}
	genesisBlock = Block{
		Index:       0,
		Timestamp:   t.String(),
		Hash:        bm.calculateHash(genesisBlock),
		PrevHash:    "",
		Transaction: "CREATE",
	}
	fmt.Printf("Genesis Block: %+v\n", genesisBlock)

	// 1. GENERATE TEST
	fmt.Printf("--------------\n")

	// "chain" and generate a new block
	fmt.Printf("Generating a new EXCHANGE block on top of genesis block\n")
	block1 := bm.GenerateBlock(genesisBlock, "EXCHANGE")
	fmt.Printf("Block 1: %+v\n", block1)

	// try to generate a new invalid block
	fmt.Printf("Generating a new EXCHANGE (same index but faulty hash)\n")

	// the fact that this block is generated 5 seconds later will result in a different hash
	t = time.Now().Local().Add(time.Second * time.Duration(5))
	faultyBlock := Block{}
	faultyBlock = Block{
		Index:       1,
		Timestamp:   t.String(),
		Hash:        bm.calculateHash(faultyBlock),
		PrevHash:    block1.PrevHash,
		Transaction: "EXCHANGE",
	}
	fmt.Printf("Faulty Block 1: %+v\n", faultyBlock)

	// 2. VALIDATION TEST
	fmt.Printf("--------------\n")
	ok := bm.isBlockValid(block1, genesisBlock)
	fmt.Printf("Is block 1 the correct successor to genesis block? (should be true): %t\n", ok)

	ok = bm.isBlockValid(faultyBlock, block1)
	fmt.Printf("Is (faulty) block 1 the correct successor to genesis block? (should be false): %t\n", ok)

}
