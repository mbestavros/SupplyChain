package server

import 
(
	"grouper"
	"blockmanager"
	"fmt"
)

type Server struct {
	gr grouper.Grouper
	bm blockmanager.Blockmanager
}

var sr Server 

//start network in case of genesis block
func (sr *Server) genesis(){
	// gr := grouper.Grouper{}
	// bm := blockmanager.Blockmanager{}
	genesisBlock := sr.bm.Genesis()
	sr.gr.StartNetwork("test","test","Test")
	fmt.Println(genesisBlock)
}

//TODO: join network if not genesis
func (sr *Server) join(){
	// gr := grouper.Grouper{}
	// bm := blockmanager.Blockmanager{}
	
	sr.gr.JoinNetwork("test", "Test", "test", "test", "Test")
	//broadcast user info to all users 
	//http 

	//get existing blockchain 
}

//TODO: send blocks 
// go routine for proof of work new block
func (sr *Server) sendBlock(block blockmanager.Block) {
	// gr := grouper.Grouper{}
	// bm := blockmanager.Blockmanager{}
	go func() {
		newBlock := sr.bm.GenerateBlock(block, "transaction info") //pass old block, new info
		fmt.Println(newBlock)
	}()

}

//TODO: receive blocks 
func (sr *Server) start(){
	// gr := grouper.Grouper{}
	// bm := blockmanager.Blockmanager{}
	for {
		newBlock := false
		//newBlock := <- bcServer.something

		if newBlock {
			//something 
		}
	}
}