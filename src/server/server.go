package server

import 
(
	"blockmanager"
	"bytes"
	"encoding/json"
	"fmt"
	"grouper"
	"net/http"
	"strconv"
	"sync"
)

type Server struct {
	gr grouper.Grouper
	bm blockmanager.Blockmanager
	bcServer []blockmanager.Block
	srv  *http.Server
}

var sr Server 

//Start network in case of genesis block
func (sr *Server) Genesis(myIp string, myPort string, myName string){
	
	genesisBlock := sr.bm.Genesis()
	sr.gr.StartNetwork(myIp, myPort, myName)
	sr.bcServer = append(sr.bcServer, genesisBlock)
}

// Join network if not genesis 
func (sr *Server) Join(friendIp string, friendPort string, myIp string, myPort string, myName string){
	
	sr.gr.JoinNetwork(friendIp, friendPort, myIp, myPort, myName)

	cli := &http.Client{}
	r, err := cli.Get("http://" + friendIp + ":" + friendPort + "/getBlock")
	defer r.Body.Close()
	var bcServer []blockmanager.Block
	var test int 
	err = json.NewDecoder(r.Body).Decode(&test)
	if err != nil {
		fmt.Println("ERROR in join in server.go:", err)
		//r is 404.
		return
	}
	fmt.Println(test)
	sr.bcServer = bcServer
}

// Generate new block and send to all peers via post request
func (sr *Server) SendBlock(block blockmanager.Block, transaction blockmanager.Transaction) {
	newBlock := sr.bm.GenerateBlock(block, transaction)
	sr.bcServer = append(sr.bcServer, newBlock)
	fmt.Println(newBlock)

	var wg sync.WaitGroup
	for _, usr := range sr.gr.Them {
		wg.Add(1)
		go func(p grouper.Peer) {
			b := new(bytes.Buffer)
			json.NewEncoder(b).Encode(newBlock)
			http.Post("http://"+p.Ip+":"+p.Port+"/verifyBlock", "application/json; charset=utf-8", b)
			wg.Done()
		}(usr)
		wg.Wait()
	}
}

// Helper for get request to get existing blockchain
func (sr *Server) helperJoinGetBlock(w http.ResponseWriter, r *http.Request){
	fmt.Println(sr.bcServer)
	json.NewEncoder(w).Encode(sr.bcServer)
	fmt.Println("encoded")
}


// Helper for receiving a block, checking if it's valid
func (sr *Server) helperVerifyBlock(w http.ResponseWriter, r *http.Request){
	newBlock := blockmanager.Block{}
	json.NewDecoder(r.Body).Decode(&newBlock)
	isValid := sr.bm.IsBlockValid(sr.bcServer[len(sr.bcServer)-1], newBlock)
	if isValid {
		sr.bcServer = append(sr.bcServer, newBlock)
		fmt.Println("appending in verify")
	}
	fmt.Println("helper verify successful")
}

// Listening on http server
func (sr *Server) start(){
	port_int, err := strconv.Atoi(sr.gr.Me.Port)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	// Adding 1 so it doesn't conflict with other server
	port := strconv.Itoa(port_int + 1)

	//sr.srv = &http.Server{Addr: ":" + port}
	serverMuxServer := http.NewServeMux()
	serverMuxServer.HandleFunc("/joinGetBlock", sr.helperJoinGetBlock)
	serverMuxServer.HandleFunc("/verifyBlock", sr.helperVerifyBlock)
	go func() {
		http.ListenAndServe(":"+port, serverMuxServer)
	}()
}