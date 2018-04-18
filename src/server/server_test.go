package server

import (
	"testing"
	"blockmanager"
	"fmt"
)

func TestIP(t *testing.T){
	sr := Server{}
	ip, _ := sr.externalIP()
	fmt.Println(ip)
}

func TestGenesis(t *testing.T) {
	t.Log("Testing genesis")
	sr := Server{}
	sr.Genesis("8080", "shreya")
	fmt.Println("finished genesis")
	sr.start()
}

func TestJoin(t *testing.T){
	t.Log("Testing joining a network")
	sr := Server{}
	sr.Genesis("8080", "shreya")
	sr.start()

	sr2 := Server{}
	
	sr2.Join("localhost", "8080", "8082", "shreya2")
	sr2.start()
	fmt.Println("sr2 them", sr2.gr.Them)
	fmt.Println("sr them", sr.gr.Them)
	
	t.Log("lol we outchea")

	t.Log("Testing sending new block")

	var transaction blockmanager.Transaction
	transaction.Type = blockmanager.Create
	transaction.OriginUser =  1
	transaction.DestinationUser = 2
	transaction.InitialTimestamp = 0
	transaction.FinalTimestamp = 0
	
	var block blockmanager.Block
	block.BlockTransaction = transaction

	sr.SendBlock(block, transaction)
	
}

func TestSendBlock(t *testing.T){
	t.Log("Testing sending new block")
	sr := Server{}

	var transaction blockmanager.Transaction
	transaction.Type = blockmanager.Create
	transaction.OriginUser =  1
	transaction.DestinationUser = 2
	transaction.InitialTimestamp = 0
	transaction.FinalTimestamp = 0
	
	var block blockmanager.Block
	block.BlockTransaction = transaction

	sr.SendBlock(block, transaction)
}

func TestVerifyBlock(t *testing.T){
	t.Log("Test verify block")

}