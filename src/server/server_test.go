package server

import (
	"testing"
	"blockmanager"
	"fmt"
)

func TestGenesis(t *testing.T) {
	t.Log("Testing genesis")
	sr := Server{}
	sr.Genesis("localhost", "8080", "shreya")
	fmt.Println("finished genesis")
	sr.start()
}

func TestJoin(t *testing.T){
	t.Log("Testing joining a network")
	sr := Server{}
	sr.Genesis("localhost", "8080", "shreya")
	sr.start()

	sr2 := Server{}
	sr2.Join("localhost", "8080", "localhost", "8081", "shreya2")
	sr2.start()
	t.Log("lol we outchea")
	
}

func TestSendBlock(t *testing.T){
	t.Log("Testing sending new block")
	sr := Server{}

	var transaction blockmanager.Transaction

	// transaction := new Transaction({
	// 	Type: Create,
	// 	OriginUser: 1,
	// 	DestinationUser: 2,
	// 	InitialTimestamp: 0,
	// 	FinalTimestamp: 0

	// });
	// block := new Block({
	// 	Index: 0,
	// 	Timestamp: 0,
	// 	Hash: "test",
	// 	PrevHash: "test",
	// 	Difficulty: 1,
	// 	Nonce: "test",
	// 	BlockTransaction: transaction
	// });
	
	var block blockmanager.Block
	sr.SendBlock(block, transaction)

}