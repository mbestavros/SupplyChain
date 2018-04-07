package main

import "time"
import "fmt"
import "os"
import "math/rand"

// this explains how to use grouper.go
func main() {
	fmt.Println("Starting demo...")
	if len(os.Args) < 2 {
		fmt.Println("Please run 'grouper #' where # is a number 1-3")
		fmt.Println("Please run all roughly the same time plus or minus a few seconds, and run 1 first")
	}
	switch os.Args[1] {
	case "1":
		fmt.Println("STARTING NETWORK W/ USER 1")
		// first, create a new grouper
		gr1 := Grouper{}
		// if you're starting a network, this is how
		// your IP, your Port, your name
		gr1.StartNetwork("localhost", "8080", "ken")
		// after some time, leave
		<-time.After(time.Duration(10+rand.Intn(10)) * time.Second)
		fmt.Println("LEAVING NETWORK W/ USER 1")
		gr1.Shutdown()
	case "2":
		fmt.Println("JOINING NETWORK W/ USER 2")
		// if you're joining a network, do this
		gr2 := Grouper{}
		// (join the network started above)
		// friend's IP, friend's port, your IP, your port, your name
		gr2.JoinNetwork("localhost", "8080", "localhost", "8081", "john")
		// after some time, leave
		<-time.After(time.Duration(10+rand.Intn(10)) * time.Second)
		fmt.Println("LEAVING NETWORK W/ USER 2")
		gr2.Shutdown()
	case "3":
		fmt.Println("JOINING NETWORK W/ USER 3")
		gr3 := Grouper{}
		gr3.JoinNetwork("localhost", "8081", "localhost", "8082", "marlon")
		// after some time, leave
		<-time.After(time.Duration(10+rand.Intn(10)) * time.Second)
		fmt.Println("LEAVING NETWORK W/ USER 3")
		gr3.Shutdown()
	default:
		fmt.Println("uhh...")
	}
}
