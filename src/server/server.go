package server

import (
	"blockmanager"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"grouper"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Server struct {
	gr       grouper.Grouper
	bm       blockmanager.Blockmanager
	bcServer []blockmanager.Block
	srv      *http.Server
}

var sr Server

func (sr *Server) externalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("are you connected to the network?")
}

//Start network in case of genesis block
func (sr *Server) Genesis(myPort string, myName string) {

	genesisBlock := sr.bm.Genesis()
	myIp, err := sr.externalIP()
	if err != nil {
		log.Debug("Error getting IP address: ", err)
	}
	sr.gr.StartNetwork(myIp, myPort, myName)
	sr.bcServer = append(sr.bcServer, genesisBlock)
}

// Join network if not genesis
func (sr *Server) Join(friendIp string, friendPort string, myPort string, myName string) {

	myIp, err := sr.externalIP()
	if err != nil {
		log.Debug("Error getting IP address: ", err)
	}
	sr.gr.JoinNetwork(friendIp, friendPort, myIp, myPort, myName)

	cli := &http.Client{}
	friendPort = increment_port(friendPort)
	r, err := cli.Get("http://" + friendIp + ":" + friendPort + "/joinGetBlock")
	defer r.Body.Close()
	var bcServer []blockmanager.Block

	err = json.NewDecoder(r.Body).Decode(&bcServer)
	if err != nil {
		log.Debug("ERROR in join in server.go:", err)
		return
	}
	sr.bcServer = bcServer
}

// Generate new block and send to all peers via post request
func (sr *Server) SendBlock(block blockmanager.Block, transaction blockmanager.Transaction) {
	newBlock := sr.bm.GenerateBlock(block, transaction)
	sr.bcServer = append(sr.bcServer, newBlock)
	fmt.Println("<< sending new block to others... >> ")

	var wg sync.WaitGroup
	for _, usr := range sr.gr.Them {
		wg.Add(1)
		go func(p grouper.Peer) {
			b := new(bytes.Buffer)
			json.NewEncoder(b).Encode(newBlock)
			port := increment_port(p.Port)
			http.Post("http://"+p.Ip+":"+port+"/verifyBlock", "application/json; charset=utf-8", b)
			wg.Done()
		}(usr)
		wg.Wait()
	}
	fmt.Println("<< block has been broadcasted! >> ")
}

func (sr *Server) NewTransaction(tran blockmanager.Transaction) {
	sr.SendBlock(sr.bcServer[len(sr.bcServer)-1], tran)
}

// Helper for get request to get existing blockchain
func (sr *Server) helperJoinGetBlock(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(sr.bcServer)
}

func (sr *Server) UndoBlock() {
	sr.bcServer = sr.bcServer[:len(sr.bcServer)-1]
}

// Helper for receiving a block, checking if it's valid
func (sr *Server) helperVerifyBlock(w http.ResponseWriter, r *http.Request) {
	newBlock := blockmanager.Block{}
	json.NewDecoder(r.Body).Decode(&newBlock)

	followsRules := sr.bm.BlockFollowsRules(newBlock, sr.bcServer)
	if followsRules {
		isValid := sr.bm.IsBlockValid(newBlock, sr.bcServer[len(sr.bcServer)-1])
		if isValid {
			sr.bcServer = append(sr.bcServer, newBlock)
		}
	}

	log.Debug("server: ", sr.gr.Me)
	log.Debug("bcserver", sr.bcServer)
}

func increment_port(old_port string) string {
	port_int, err := strconv.Atoi(old_port)
	if err != nil {
		log.Debug("Error: ", err)
		return "ERROR"
	}

	// Adding 1 so it doesn't conflict with other server
	new_port := strconv.Itoa(port_int + 1)
	return new_port
}

// Listening on http server
func (sr *Server) Start() {
	port := increment_port(sr.gr.Me.Port)
	serverMuxServer := http.NewServeMux()
	sr.srv = &http.Server{Handler: serverMuxServer}
	serverMuxServer.HandleFunc("/joinGetBlock", sr.helperJoinGetBlock)
	serverMuxServer.HandleFunc("/verifyBlock", sr.helperVerifyBlock)
	go func() {
		http.ListenAndServe(":"+port, serverMuxServer)
	}()
}

func (sr *Server) Shutdown() {
	if sr.srv != nil {
		fmt.Println("<< shutting down server... >>")
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		sr.gr.Shutdown()
		sr.srv.Shutdown(ctx)
	}
}
