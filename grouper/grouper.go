package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
)

type peer struct {
	Name string
	Ip   string
	Port string
}

type Grouper struct {
	mu   sync.Mutex
	Me   peer
	Them []peer
	srv  *http.Server
}

var gr Grouper

// USAGE
// run 'grouper start ip port name' to begin your own network
// (params are your own ip, port, and name)
// run 'grouper join ipFriend portFriend ipSelf portSelf name' to join a network
// (ipFriend is the node alread in the network, self and name
// for this actual node)

/*
func main() {
	cmd := os.Args[1]
	gr = Grouper{}

	// handle a start, join, or help command
	if cmd == "start" {
		gr.StartNetwork(os.Args[2], os.Args[3], os.Args[4])
	} else if cmd == "join" {
		gr.JoinNetwork(os.Args[2], os.Args[3], os.Args[4], os.Args[5], os.Args[6])
	} else if cmd == "help" {
		fmt.Println("To start a network:")
		fmt.Println("\tgrouper start ip_addr port_number name")
		fmt.Println("To join a network:")
		fmt.Println("\tgrouper join friend_ip friend_port my_ip my_port name")
	} else {
		fmt.Println("Unrecognized command:", cmd)
	}

	// wait (forever) for a SIGINT
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
*/

// func - start your own network
func (gr *Grouper) StartNetwork(myIp, myPort, myName string) {
	go gr.listenToSIGINT()
	fmt.Println("starting network as", myIp+":"+myPort+":"+myName)
	gr.Me = peer{Ip: myIp, Port: myPort, Name: myName}
	gr.startHttpServer()
}

// func - join a network
func (gr *Grouper) JoinNetwork(friendIp, friendPort, myIp, myPort, myName string) {
	go gr.listenToSIGINT()
	fmt.Println("joining network as", myIp+":"+myPort+":"+myName)
	gr.Me = peer{Ip: myIp, Port: myPort, Name: myName}
	gr.getPeers(peer{Ip: friendIp, Port: friendPort})
	gr.sendJoinRequests()
	gr.startHttpServer()
}

func (gr *Grouper) getPeers(friend peer) {
	cli := &http.Client{}
	r, err := cli.Get("http://" + friend.Ip + ":" + friend.Port + "/getPeers")
	if err != nil {
		return
	}
	defer r.Body.Close()
	var otherUsers []peer
	err = json.NewDecoder(r.Body).Decode(&otherUsers)
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}
	gr.Them = otherUsers
}

func (gr *Grouper) sendJoinRequests() {
	for _, usr := range gr.Them {
		go func(p peer) {
			b := new(bytes.Buffer)
			json.NewEncoder(b).Encode(gr.Me)
			http.Post("http://"+p.Ip+":"+p.Port+"/joinNet", "application/json; charset=utf-8", b)
		}(usr)
	}
}

func (gr *Grouper) sendLeaveRequests() {
	var wg sync.WaitGroup
	for _, usr := range gr.Them {
		wg.Add(1)
		go func(p peer) {
			b := new(bytes.Buffer)
			json.NewEncoder(b).Encode(gr.Me)
			http.Post("http://"+p.Ip+":"+p.Port+"/leaveNet", "application/json; charset=utf-8", b)
			wg.Done()
		}(usr)
		wg.Wait()
	}
}

func (gr *Grouper) listenToSIGINT() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		gr.Shutdown()
		os.Exit(0)
	}()
}

func (gr *Grouper) Shutdown() {
	gr.sendLeaveRequests()
	gr.srv.Shutdown(nil)
}

func (gr *Grouper) startHttpServer() {
	gr.srv = &http.Server{Addr: ":" + gr.Me.Port}
	http.HandleFunc("/getPeers", gr.handleGetPeers)
	http.HandleFunc("/joinNet", gr.handleJoinNet)
	http.HandleFunc("/leaveNet", gr.handleLeaveNet)
	go func() {
		http.ListenAndServe(":"+gr.Me.Port, nil)
	}()
}

func (gr *Grouper) handleGetPeers(w http.ResponseWriter, r *http.Request) {
	allUsers := append(gr.Them, gr.Me)
	json.NewEncoder(w).Encode(allUsers)
}

func (gr *Grouper) handleJoinNet(w http.ResponseWriter, r *http.Request) {
	user := peer{}
	json.NewDecoder(r.Body).Decode(&user)
	gr.Them = append(gr.Them, user)
	fmt.Println(user.Name, "has joined")
}

func (gr *Grouper) handleLeaveNet(w http.ResponseWriter, r *http.Request) {
	user := peer{}
	json.NewDecoder(r.Body).Decode(&user)
	fmt.Println(user.Name, "has left")
	// find who it is, and remove them
	for ind, usr := range gr.Them {
		if user == usr {
			gr.Them = append(gr.Them[0:ind], gr.Them[ind+1:len(gr.Them)]...)
			break
		}
	}
}
