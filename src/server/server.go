package server

import 
(
	"grouper"
	"blockmanager"
)

type Server struct {
	//http probably
}

var sr Server 

//start network in case of genesis block
func (sr *Server) genesis(){
	gr := grouper.Grouper{}
	bm := blockmanager.Blockmanager{}
	genesisBlock := bm.Genesis()
	gr.StartNetwork("test","test","Test")
}

//TODO: join network if not genesis
func (sr *Server) join(){
	gr := grouper.Grouper{}
	bm := blockmanager.Blockmanager{}
	
	gr.JoinNetwork("test", "Test", "test", "test", "Test")
	//broadcast user info to all users 
	//http 

	//get existing blockchain 
}

//TODO: send blocks 
// go routine for proof of work new block
func (sr *Server) sendBlock(block Block) {
	gr := grouper.Grouper{}
	bm := blockmanager.Blockmanager{}
	go func() {
		newBlock := bm.GenerateBlock(block, "transaction info") //pass old block, new info
	}()

}

//TODO: receive blocks 
func (sr *Server) start(){
	gr := grouper.Grouper{}
	bm := blockmanager.Blockmanager{}
	for {
		newBlock := false
		//newBlock := <- bcServer.something

		if newBlock {
			//something 
		}
	}
}