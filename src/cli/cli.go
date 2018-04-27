package main

import (
	"blockmanager"
	"bufio"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"server"
	"strings"
	// "time"
)

type Cli struct {
	bm      blockmanager.Blockmanager
	sr      server.Server
	started bool
	myName  string
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
		case "undo block":
			cl.sr.UndoBlock()
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
		cl.myName = myName
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
		cl.myName = myName
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
		fmt.Println("Creating an item, adding it to the blockchain...")
		tran := cl.bm.BuildCreateTransaction(itemName, cl.myName)
		fmt.Println("ItemID", tran.Cr.ItemId)
		cl.sr.NewTransaction(tran)
	case "exchange":
		itemName := readString("What is the item? ")
		itemId := readString("What is the ID? ")
		recipientName := readString("Who are you sending it to? ")
		fmt.Println("Exchanging it on the blockchain...")
		tran := cl.bm.BuildExchangeTransaction(itemName, itemId, cl.myName, recipientName)
		cl.sr.NewTransaction(tran)
	case "consume":
		itemId := readString("What is the ID? ")
		fmt.Println("Consuming it from the blockchain...")
		tran := cl.bm.BuildConsumeTransaction(itemId, cl.myName)
		cl.sr.NewTransaction(tran)
	case "make":
		itemNames := strings.Split(readString("What are the items? (List names separated by commas) "), ",")
		itemIDs := strings.Split(readString("What are the items? (List IDs separated by commas) "), ",")
		outputItem := readString("What are you making? ")
		fmt.Println("Making it on the blockchain...")
		tran := cl.bm.BuildMakeTransaction(itemNames, itemIDs, outputItem, cl.myName)
		fmt.Println(tran.Ma.OutputItemName, ":", tran.Ma.OutputItemId)
		cl.sr.NewTransaction(tran)
	case "split":
		inputItemName := readString("What are you splitting? (name) ")
		inputItemID := readString("What are you splitting? (ID) ")
		outputNames := strings.Split(readString("What are the items? (List names separated by commas) "), ",")
		fmt.Println("Splitting it on the blockchain...")
		recipients := make([]string, len(outputNames))
		for i := range recipients {
			recipients[i] = cl.myName
		}
		tran := cl.bm.BuildSplitTransaction(inputItemName, inputItemID, outputNames, cl.myName, recipients)
		for ind, name := range tran.Sp.OutputItemNames {
			fmt.Println(name, ":", tran.Sp.OutputItemIds[ind])
		}
		cl.sr.NewTransaction(tran)
	default:
		fmt.Println("Transaction type not recognized.")
	}
}

func (cl *Cli) lookupFunc() {
	itemID := readString("Item ID: ")
	cl.sr.LookupItem(itemID)
}

func (cl *Cli) quitCommand() {
	if cl.started {
		cl.sr.Shutdown()
	}
}
