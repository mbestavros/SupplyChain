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
	BlockTransaction TransactionProvider
}

type Blockmanager struct {
}

func main() {
	fmt.Println("hello?")
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
func (bm *Blockmanager) GenerateBlock(oldBlock Block, transaction TransactionProvider) Block {
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
	t1 := &CreateTransaction{}
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

//Returns a new Create Transaction struct
func (bm *Blockmanager) BuildCreateTransaction(itemName string, userId string) TransactionProvider {
	currTime := int64(time.Now().Unix())

	// origin user is empty string because an item is being created
	createTrans := &CreateTransaction{
		Transaction: Transaction{
			TransactionType: Create,
			TimeTransacted:  currTime,
		},
		OriginUserId:      "",
		DestinationUserId: userId,
		ItemName:          itemName,
		ItemId:            generateUID(),
	}

	return createTrans
}

//Returns a new Exchange Transaction struct
func (bm *Blockmanager) BuildExchangeTransaction(itemName string, originUserId string, destinationUserId string) TransactionProvider {
	currTime := int64(time.Now().Unix())
	exchTrans := &ExchangeTransaction{
		Transaction: Transaction{
			TransactionType: Exchange,
			TimeTransacted:  currTime,
		},
		OriginUserId:      originUserId,
		DestinationUserId: destinationUserId,
		ItemName:          itemName,
		ItemId:            generateUID(),
	}

	return exchTrans
}

//Returns a new Consume Transaction struct
func (bm *Blockmanager) BuildConsumeTransaction(itemName string, consumerUserId string) TransactionProvider {
	currTime := int64(time.Now().Unix())

	// originUser == destinationUser because destination is consumer themselves
	consumeTrans := &ConsumeTransaction{
		Transaction: Transaction{
			TransactionType: Consume,
			TimeTransacted:  currTime,
		},
		OriginUserId:      consumerUserId,
		DestinationUserId: consumerUserId,
		ItemName:          itemName,
		ItemId:            generateUID(),
	}

	return consumeTrans
}

//Returns a new Make Transaction struct
func (bm *Blockmanager) BuildMakeTransaction(inputItemNames []string, inputItemIds []string, outputItemName string, makerUserId string) TransactionProvider {
	currTime := int64(time.Now().Unix())

	// originUser == destinationUser because destination is maker themselves
	// list of input items -> one output item
	makeTrans := &MakeTransaction{
		Transaction: Transaction{
			TransactionType: Make,
			TimeTransacted:  currTime,
		},
		OriginUserId:      makerUserId,
		DestinationUserId: makerUserId,
		InputItemIds:      inputItemIds,
		InputItemNames:    inputItemNames,
		OutputItemName:    outputItemName,
		OutputItemId:      generateUID(),
	}

	return makeTrans
}

//Returns a new Split Transaction struct
func (bm *Blockmanager) BuildSplitTransaction(inputItemName string, inputItemId string, outputItemNames []string, originUserId string, destinationUserIds []string) TransactionProvider {
	currTime := int64(time.Now().Unix())

	// generate new item Id's for each new item that's been split
	var outputItemIds []string
	for range outputItemNames {
		outputItemIds = append(outputItemIds, generateUID())
	}

	splitTrans := &SplitTransaction{
		Transaction: Transaction{
			TransactionType: Split,
			TimeTransacted:  currTime,
		},
		OriginUserId:       originUserId,
		DestinationUserIds: destinationUserIds,
		InputItemId:        inputItemId,
		InputItemName:      inputItemName,
		OutputItemNames:    outputItemNames,
		OutputItemIds:      outputItemIds,
	}

	return splitTrans
}
