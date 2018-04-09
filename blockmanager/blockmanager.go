package main

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"time"
)

// Block represents each 'item' in the blockchain
type Block struct {
	Index     int
	Timestamp string
	Hash      string
	PrevHash  string

	// other transaction properties within a block
	Transaction string
}

type Blockmanager struct {
}

var bm Blockmanager

// make sure block is valid by checking index, and comparing the hash of the previous block
func (bm *Blockmanager) isBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	if bm.calculateHash(newBlock) != newBlock.Hash {
		return false
	}

	return true
}

// SHA256 hasing
func (bm *Blockmanager) calculateHash(block Block) string {
	record := strconv.Itoa(block.Index) + block.Timestamp + block.Transaction + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// create a new block using previous block's hash
func (bm *Blockmanager) generateBlock(oldBlock Block, transaction string) Block {

	var newBlock Block

	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.Transaction = transaction
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = bm.calculateHash(newBlock)

	return newBlock
}
