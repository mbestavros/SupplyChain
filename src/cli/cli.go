package main

import (
	"blockmanager"
	"bufio"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"server"
	"time"
)

type Cli struct {
	sr      server.Server
	started bool
}

func readString(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	if prompt != "" {
		fmt.Print(prompt)
	}
	entered, _ := reader.ReadString('\n')
	entered = entered[:len(entered)-1]
	return entered
}

// this main function is the main function of the entire program.
// It all coordinates here.
func main() {
	if len(os.Args) > 1 && os.Args[1] == "debug" {
		log.SetLevel(log.DebugLevel)
		log.Debug("Debug statement!")
	}
	cl := Cli{}
	fmt.Println()
	entered := "welcome"
	for entered != "quit" {
		switch entered {
		case "welcome":
			fmt.Println("Welcome to our CS451 / CS651 final project!")
			fmt.Println("To view a list of commands, type \"help\".")
		case "help":
			fmt.Println("Help Menu:")
			fmt.Println(" help - view this help menu \n quit - exit the program \n start - start a brand new blockchain")
			fmt.Println(" join - join an existing blockchain network \n transact - create a transaction")
			fmt.Println(" lookup - look up the status of an item ")
		case "start":
			cl.startFunc()
		case "join":
			cl.joinFunc()
		case "transact":
			cl.transactFunc()
		case "lookup":
			cl.lookupFunc()
		default:
			fmt.Println("Unrecognized command. Type \"help\" for a list of commands")
		}
		entered = readString("> ")
		fmt.Println()
		if entered == "quit" {
			cl.quitCommand()
		}
	}
	fmt.Println("Goodbye.")
}

func (cl *Cli) startFunc() {
	if cl.started {
		fmt.Println("Cant start or join twice. Quit and try again")
	} else {
		myPort := readString("Your port: ")
		myName := readString("Your name: ")
		cl.sr = server.Server{}
		cl.sr.Genesis(myPort, myName)
		cl.sr.Start()
		cl.started = true
	}
}

func (cl *Cli) joinFunc() {
	if cl.started {
		fmt.Println("Cant start or join twice. Quit and try again")
	} else {
		friendIP := readString("Your friend's IP address: ")
		friendPort := readString("Your friend's port: ")
		myPort := readString("Your port: ")
		myName := readString("Your name: ")
		cl.sr = server.Server{}
		cl.sr.Join(friendIP, friendPort, myPort, myName)
		cl.sr.Start()
		cl.started = true
	}
}

func (cl *Cli) transactFunc() {
	transactType := readString("What kind of transaction?\n(create, exchange, consume, make, split) ")
	switch transactType {
	case "create":
		itemName := readString("What is the item? ")
		// TODO: create the item's transaction, make a block, verify it, and send it out
		fmt.Println(itemName, "has an ID of xxxxxxx")
		fmt.Println("Now to mine the relevant block and add it to the blockchain...")
		tran := blockmanager.Transaction{
			Type:             "create",
			OriginUser:       0,
			DestinationUser:  0,
			InitialTimestamp: time.Now().Second(),
			Hash:             "temporary"}
		cl.sr.NewTransaction(tran)
	case "exchange":
	case "consume":
	case "make":
	case "split":
	default:
		fmt.Println("Transaction type not recognized.")
	}
}

func (cl *Cli) lookupFunc() {
	itemID := readString("Item ID: ")
	fmt.Println("Completed lookup for", itemID)
}

func (cl *Cli) quitCommand() {
	if cl.started {
		cl.sr.Shutdown()
	}
}
