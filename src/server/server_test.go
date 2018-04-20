package server

import (
	"blockmanager"
	"fmt"
	"testing"
	"time"
)

func TestIP(t *testing.T) {
	sr := Server{}
	ip, _ := sr.externalIP()
	fmt.Println(ip)
}

func TestJoin(t *testing.T) {
	t.Log("Testing genesis block")
	fmt.Println("Testing genesis block")

	sr := Server{}
	sr.Genesis("8080", "shreya")
	sr.start()

	// sleep so that sr has time to start up
	time.Sleep(1000 * time.Millisecond)
	t.Log("Testing joining a network")

	sr2 := Server{}
	sr2.Join("localhost", "8080", "8082", "shreya2")
	sr2.start()

	// sleep so that sr2 and sr are all caught up before printing
	time.Sleep(1000 * time.Millisecond)

	t.Log("Testing sending new block")
	fmt.Println("Testing sending and verifying new block")

	var transaction blockmanager.Transaction
	transaction.Type = blockmanager.Create
	transaction.OriginUser = 1
	transaction.DestinationUser = 2
	transaction.InitialTimestamp = 0
	transaction.FinalTimestamp = 0


	sr.SendBlock(sr.bcServer[0], transaction)

	time.Sleep(1000 * time.Millisecond)
	if len(sr.bcServer) == len(sr2.bcServer) {
		fmt.Println("Blockchain is the same for servers 1 and 2")
	} else {
		return
	}
	sr.shutdown()
	sr2.shutdown()
}
