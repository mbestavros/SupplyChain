package blockmanager

import (
	"fmt"
	"testing"
	"time"
)

func TestGenerateMiningVerify(t *testing.T) {
	// initialize block manager
	bm := Blockmanager{}

	fmt.Println("Generating a genesis block (CREATE transaction)...")

	currTime := time.Now()
	genesisBlock := bm.genesis()
	fmt.Printf("Genesis Block: %+v\n", genesisBlock)

	// 1. GENERATE TEST
	fmt.Printf("--------------\n")

	// "chain" and generate a new block
	t2 := Transaction{
		Type: "Create",
	}
	fmt.Printf("Generating a new EXCHANGE block on top of genesis block\n")
	block1 := bm.GenerateBlock(genesisBlock, t2)
	fmt.Printf("Block 1: %+v\n", block1)

	// try to generate a new invalid block
	fmt.Printf("Generating a new EXCHANGE (same index but faulty hash)\n")

	// the fact that this block is generated 5 seconds later will result in a different hash
	currTime = time.Now().Local().Add(time.Second * time.Duration(5))
	faultyBlock := Block{}
	faultyBlock = Block{
		Index:            1,
		Timestamp:        currTime.String(),
		Hash:             block1.Hash,
		PrevHash:         genesisBlock.Hash,
		BlockTransaction: t2,
	}
	fmt.Printf("Faulty Block 1: %+v\n", faultyBlock)

	// 2. VALIDATION TEST
	fmt.Printf("--------------\n")
	ok := bm.IsBlockValid(block1, genesisBlock)
	fmt.Printf("Is block 1 the correct successor to genesis block? (should be true): %t\n", ok)

	ok = bm.IsBlockValid(faultyBlock, genesisBlock)
	fmt.Printf("Is (faulty) block 1 the correct successor to genesis block? (should be false): %t\n", ok)

}
