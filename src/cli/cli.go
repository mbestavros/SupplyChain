package main

import (
	"bufio"
	"fmt"
	"grouper"
	"os"
)

type Cli struct {
	gr grouper.Grouper
}

// this main function is the main function of the entire program.
// It all coordinates here.
func main() {
	fmt.Println()
	reader := bufio.NewReader(os.Stdin)
	entered := "welcome"
	for entered != "quit" {
		switch entered {
		case "welcome":
			fmt.Println("Welcome to our CS451 / CS651 final project!")
			fmt.Println("To view a list of commands, type \"help\"")
		case "help":
			fmt.Println("Help Menu:")
			fmt.Println(" help - view this help menu \n quit - exit the program \n start - start a brand new blockchain")
			fmt.Println(" join - join an existing blockchain network \n transact - create a transaction")
		default:
			fmt.Println("Unrecognized command. Type \"help\" for a list of commands")
		}
		fmt.Print("> ")
		entered, _ = reader.ReadString('\n')
		entered = entered[:len(entered)-1]
		fmt.Println()
	}
	fmt.Println("Goodbye.")
}
