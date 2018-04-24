package blockmanager

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/rs/xid"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

const miningDifficulty = 1

// Block represents each 'item' in the blockchain
type Block struct {
	Index     int
	Timestamp string
	Hash      string
	PrevHash  string

	// attributes used for mining
	Difficulty int
	Nonce      string

	// other transaction properties within a block
	BlockTransaction Transaction
}

type Blockmanager struct {
}

type Action string

func main() {
	fmt.Println("hello?")
}

//Enum for all types of transactions
const (
	Create   Action = "Create"
	Exchange Action = "Exchange"
	Consume  Action = "Consume"
	Make     Action = "Make"
	Split    Action = "Split"
)

//Struct to represent all Transactions
type Transaction struct {
	Type             Action
	OriginUser       int64
	DestinationUser  int64
	InitialTimestamp int
	FinalTimestamp   int

	//Hash probably goes here as well @Sean
	Hash string
}

var bm Blockmanager

func generateUID() string {
	return xid.New().String()
}

// make sure block is valid by checking index, and comparing the hash of the previous block
func (bm *Blockmanager) IsBlockValid(newBlock Block, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	if bm.calculateHash(newBlock) != newBlock.Hash {
		return false
	}

	return true
}

// SHA256 hasing
func (bm *Blockmanager) calculateHash(block Block) string {
	// convert transaction struct into json string to be hashed
	b, err := json.Marshal(block.BlockTransaction)
	if err != nil {
		fmt.Println(err)
	}
	transactionString := string(b)

	record := strconv.Itoa(block.Index) + block.Timestamp + transactionString + block.PrevHash + block.Nonce
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// create a new block using previous block's hash
func (bm *Blockmanager) GenerateBlock(oldBlock Block, transaction Transaction) Block {
	var newBlock Block

	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.BlockTransaction = transaction
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Difficulty = miningDifficulty

	for i := 0; ; i++ {
		hex := fmt.Sprintf("%x", i)
		newBlock.Nonce = hex
		fmt.Println("<< mining... >>")
		hashAttempt := bm.calculateHash(newBlock)
		if !bm.isHashValid(hashAttempt, newBlock.Difficulty) {
			log.Debug(hashAttempt, " not valid. Trying again...")
			// simulate proof of work time consumed
			time.Sleep(time.Second)
			continue
		} else {
			log.Debug(hashAttempt, " valid! Block mined")
			fmt.Println("<< mined! >>")
			newBlock.Hash = hashAttempt
			break
		}

	}

	return newBlock
}

func (bm *Blockmanager) Genesis() Block {
	t1 := Transaction{
		Type: "Create",
	}
	currTime := time.Now()
	genesisBlock := Block{}
	genesisBlock = Block{
		Index:            0,
		Timestamp:        currTime.String(),
		Hash:             bm.calculateHash(genesisBlock),
		PrevHash:         "",
		BlockTransaction: t1,
	}
	//spew.Dump(genesisBlock)
	return genesisBlock

}

func (bm *Blockmanager) isHashValid(hash string, difficulty int) bool {
	// temporarily making hash easier (avg 5 tries)
	prefix := strings.Repeat("0", difficulty)
	prefix2 := strings.Repeat("1", difficulty)
	prefix3 := strings.Repeat("2", difficulty)
	return strings.HasPrefix(hash, prefix) ||
		strings.HasPrefix(hash, prefix2) ||
		strings.HasPrefix(hash, prefix3)

}
