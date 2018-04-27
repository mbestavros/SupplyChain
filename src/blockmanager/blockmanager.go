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
		fmt.Println(" << recieved invalid block: index mismatch >>")
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		fmt.Println(" << recieved invalid block: previous hash mismatch >>")
		return false
	}

	if bm.calculateHash(newBlock) != newBlock.Hash {
		fmt.Println(" << recieved invalid block: self hash mismatch >>")
		return false
	}

	fmt.Println(" << recieved valid block:", newBlock.BlockTransaction.TransactionType, " >>")
	// more logging here
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
	t1 := Transaction{Cr: CreateTransaction{}}
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
func (bm *Blockmanager) BuildCreateTransaction(itemName string, userId string) Transaction {
	return Transaction{
		TransactionType: Create,
		TimeTransacted:  int64(time.Now().Unix()),
		Cr: CreateTransaction{
			OriginUserId:      "",
			DestinationUserId: userId,
			ItemName:          itemName,
			ItemId:            generateUID(),
		},
	}

}

//Returns a new Exchange Transaction struct
func (bm *Blockmanager) BuildExchangeTransaction(itemName string, originUserId string, destinationUserId string) Transaction {
	return Transaction{
		TransactionType: Exchange,
		TimeTransacted:  int64(time.Now().Unix()),
		Ex: ExchangeTransaction{
			OriginUserId:      originUserId,
			DestinationUserId: destinationUserId,
			ItemName:          itemName,
			ItemId:            generateUID(),
		},
	}
}

//Returns a new Consume Transaction struct
func (bm *Blockmanager) BuildConsumeTransaction(itemName string, consumerUserId string) Transaction {
	return Transaction{
		TransactionType: Consume,
		TimeTransacted:  int64(time.Now().Unix()),
		Co: ConsumeTransaction{
			OriginUserId:      consumerUserId,
			DestinationUserId: consumerUserId,
			ItemName:          itemName,
			ItemId:            generateUID(),
		},
	}

}

//Returns a new Make Transaction struct
func (bm *Blockmanager) BuildMakeTransaction(inputItemNames []string, inputItemIds []string, outputItemName string, makerUserId string) Transaction {
	return Transaction{
		TransactionType: Make,
		TimeTransacted:  int64(time.Now().Unix()),
		Ma: MakeTransaction{
			OriginUserId:      makerUserId,
			DestinationUserId: makerUserId,
			InputItemIds:      inputItemIds,
			InputItemNames:    inputItemNames,
			OutputItemName:    outputItemName,
			OutputItemId:      generateUID(),
		},
	}
}

//Returns a new Split Transaction struct
func (bm *Blockmanager) BuildSplitTransaction(inputItemName string, inputItemId string, outputItemNames []string, originUserId string, destinationUserIds []string) Transaction {
	// generate new item Id's for each new item that's been split
	var outputItemIds []string
	for range outputItemNames {
		outputItemIds = append(outputItemIds, generateUID())
	}

	return Transaction{
		TransactionType: Split,
		TimeTransacted:  int64(time.Now().Unix()),
		Sp: SplitTransaction{
			OriginUserId:       originUserId,
			DestinationUserIds: destinationUserIds,
			InputItemId:        inputItemId,
			InputItemName:      inputItemName,
			OutputItemNames:    outputItemNames,
			OutputItemIds:      outputItemIds,
		},
	}

}

// History of an item
func (bm *Blockmanager) GetItemHistory(itemId string, bcServer []Block) []Transaction{
	var result []Transaction

	// Start at 1 to skip genesis
	for block_i := 1; block_i < len(bcServer); block_i++ {
		block := bcServer[block_i]
		transaction := block.BlockTransaction

		switch transType := transaction.TransactionType; transType {
		
		case Split:
			splitTrans := transaction.Sp
			if splitTrans.InputItemId == itemId {

				result = append(result, transaction)
			} else {
				for i := range splitTrans.OutputItemIds {
					outputId := splitTrans.OutputItemIds[i]
					if outputId == itemId {
						result = append(result, transaction)
					}
				}
			}
		
		case Make:
			makeTrans := transaction.Ma
			if makeTrans.OutputItemId == itemId {
				result = append(result, transaction)
			} else {
				for i := range makeTrans.InputItemIds {
					inputId := makeTrans.InputItemIds[i]
					if inputId == itemId {
						result = append(result, transaction)
					}
				}
			}
			
		case Create:
			createTrans := transaction.Cr
			if createTrans.ItemId == itemId {
				result = append(result, transaction)
			}

		case Exchange:
			exchangeTrans := transaction.Ex
			if exchangeTrans.ItemId == itemId {
				result = append(result, transaction)
			}

		case Consume:
			consumeTrans := transaction.Co
			if consumeTrans.ItemId == itemId {
				result = append(result, transaction)
			}

		default:
			fmt.Println("error: none of the transaction types have a value")
			return result
		}
	}
	return result
}
