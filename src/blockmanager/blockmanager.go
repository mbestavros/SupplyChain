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
		fmt.Println(" << invalid block: index mismatch >>")
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		fmt.Println(" << invalid block: previous hash mismatch >>")
		return false
	}

	if bm.calculateHash(newBlock) != newBlock.Hash {
		fmt.Println(" << invalid block: self hash mismatch >>")
		return false
	}

	// VALID BLOCK
	fmt.Println("<< valid block:", newBlock.BlockTransaction.TransactionType, " >>")
	trans := newBlock.BlockTransaction
	switch trans.TransactionType {
	case Create:
		fmt.Println("<<", trans.Cr.OriginUserId, "created", trans.Cr.ItemName, ">>")
	case Exchange:
		fmt.Println("<<", trans.Ex.OriginUserId, "exchanged",
			trans.Ex.ItemName, "to", trans.Ex.DestinationUserId, ">>")
	case Consume:
		fmt.Println("<<", trans.Co.OriginUserId, "consumed", trans.Co.ItemName, ">>")
	case Make:
		fmt.Println("<<", trans.Ma.OriginUserId, "maked", trans.Ma.OutputItemName, "from", trans.Ma.InputItemNames, ">>")
	case Split:
		fmt.Println("<<", trans.Sp.OriginUserId, "split", trans.Sp.InputItemName, "into", trans.Sp.OutputItemNames, ">>")
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
	// NOTE TEMPORARILY ALWAYS TRUE
	return true || strings.HasPrefix(hash, prefix) ||
		strings.HasPrefix(hash, prefix2) ||
		strings.HasPrefix(hash, prefix3)
}

//Returns a new Create Transaction struct
func (bm *Blockmanager) BuildCreateTransaction(itemName, userId string) Transaction {
	return Transaction{
		TransactionType: Create,
		TimeTransacted:  int64(time.Now().Unix()),
		Cr: CreateTransaction{
			OriginUserId:      userId,
			DestinationUserId: userId,
			ItemName:          itemName,
			ItemId:            generateUID(),
		},
	}

}

//Returns a new Exchange Transaction struct
func (bm *Blockmanager) BuildExchangeTransaction(itemName, itemId, originUserId, destinationUserId string) Transaction {
	return Transaction{
		TransactionType: Exchange,
		TimeTransacted:  int64(time.Now().Unix()),
		Ex: ExchangeTransaction{
			OriginUserId:      originUserId,
			DestinationUserId: destinationUserId,
			ItemName:          itemName,
			ItemId:            itemId,
		},
	}
}

//Returns a new Consume Transaction struct
func (bm *Blockmanager) BuildConsumeTransaction(itemId, consumerUserId string) Transaction {
	return Transaction{
		TransactionType: Consume,
		TimeTransacted:  int64(time.Now().Unix()),
		Co: ConsumeTransaction{
			OriginUserId:      consumerUserId,
			DestinationUserId: consumerUserId,
			ItemName:          "",
			ItemId:            itemId,
		},
	}

}

//Returns a new Make Transaction struct
func (bm *Blockmanager) BuildMakeTransaction(inputItemNames, inputItemIds []string, outputItemName, makerUserId string) Transaction {
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
func (bm *Blockmanager) BuildSplitTransaction(inputItemName, inputItemId string, outputItemNames []string, originUserId string, destinationUserIds []string) Transaction {
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

func (tr *Transaction) TransOwner() string {
	switch tr.TransactionType {
	case Create:
		return tr.Cr.DestinationUserId
	case Exchange:
		return tr.Ex.DestinationUserId
	case Consume:
		return ""
	case Make:
		return tr.Ma.DestinationUserId
	case Split:
		return tr.Sp.DestinationUserIds[0]
	}
	return ""
}

// handles item stealing, double spending, ...
// this isn't exhaustive, but demonstrates the point of the blockchain
func (bm *Blockmanager) BlockFollowsRules(block Block, bc []Block) bool {
	trans := block.BlockTransaction
	switch trans.TransactionType {
	case Create:
		// TODO
	case Exchange:
		// make sure the sender actually owns it
		// make sure it hasn't been consumed
		// (no double spending, no stealing)
		hist := bm.GetItemHistory(block.BlockTransaction.Ex.ItemId, bc)
		if hist == nil || len(hist) == 0 {
			fmt.Println("<< invalid block: item has no history >>")
			return false
		}
		lastTrans := hist[len(hist)-1]
		switch lastTrans.TransactionType {
		case Create:
			if lastTrans.TransOwner() != block.BlockTransaction.Ex.OriginUserId {
				fmt.Println("<< invalid block: incorrect owner >>")
				// fmt.Println("last owner:", lastTrans.TransOwner())
				return false
			}
		case Exchange:
			if lastTrans.TransOwner() != block.BlockTransaction.Ex.OriginUserId {
				fmt.Println("<< invalid block: incorrect owner >>")
				return false
			}
		case Consume:
			fmt.Println("<< invalid block: item has been consumed >>")
			return false
		case Make:
			if lastTrans.TransOwner() != block.BlockTransaction.Ex.OriginUserId {
				fmt.Println("<< invalid block: incorrect owner >>")
				return false
			}
		case Split:
			if lastTrans.TransOwner() != block.BlockTransaction.Ex.OriginUserId {
				fmt.Println("<< invalid block: incorrect owner >>")
				return false
			}
		}
	case Consume:
		// TODO
	case Make:
		// TODO
	case Split:
		// TODO
	}
	return true
}

//A function that takes a userID as an input and returns the Items they currently own.
func (bm *Blockmanager) GetItemsOfOwner(userID string, bcServer []Block) []Transaction{
	var result []Transaction
	transactionsMap := make(map[string]string)
	for block_i := 1; block_i < len(bcServer); block_i++ {
		block := bcServer[block_i]
		transaction := block.BlockTransaction
		switch transType := transaction.TransactionType; transType {

		//If the destinationUser of the block matches the userID, we check to see if that block is already in the map.
		//If it is, we delete it and replace it with this transaction, to ensure we don't tell the user they own a block
		//that they had previously sold. Otherwise, we just add it to the map.
		//Note that they values could really be anything, in the end I just print out the keys.

		case Split:
			splitTrans := transaction.Sp
			if splitTrans.OriginUserId == userID{ // If the origin is the user, then they must no longer own the block
				delete(transactionsMap, splitTrans.OutputItemName)
			} else {
			for i := 0; i <= len(splitTrans.DestinationUserIds); i++{
				if splitTrans.DestinationUserIds[i] == userID {
					if val, ok := transactionsMap[splitTrans.OutputItemName]; ok{
						delete(transactionsMap, splitTrans.OutputItemName)
					}
					transactionsMap[splitTrans.OutputItemName] = userID
					break
				}
			}
		}


		case Make:
			makeTrans := transaction.Ma
			if makeTrans.DestinationUserId == userID {
				if val, ok := transactionsMap[makeTrans.OutputItemName]; ok{
					delete(transactionsMap, makeTrans.OutputItemName)
				}
				transactionsMap[makeTrans.OutputItemName] = userID
			}

		case Create:
			createTrans := transaction.Cr
			if createTrans.DestinationUserId == userID {
				if val, ok := transactionsMap[createTrans.ItemName]; ok{
					delete(transactionsMap, makeTrans.ItemName)
				}
				transactionsMap[makeTrans.ItemName] = userID
			}

		case Exchange:
			exchangeTrans := transaction.Ex
			if exchangeTrans.OriginUserId == userID{ // If the origin is the user, then they must no longer own the block
				delete(transactionsMap, exchangeTrans.ItemName)
			} else if exchangeTrans.DestinationUserId == userID {
				if val, ok := transactionsMap[exchangeTrans.ItemName]; ok{
					delete(transactionsMap, exchangeTrans.ItemName)
				}
				transactionsMap[makeTrans.ItemName] = userID
			}

		case Consume:
			consumeTrans := transaction.Co
			if consumeTrans.DestinationUserId == userID {
				if val, ok := transactionsMap[consumeTrans.ItemName]; ok{
					delete(transactionsMap, consumeTrans.ItemName)
				}
				transactionsMap[consumeTrans.ItemName] = userID
			}
	}
	var result []Transaction
	for k := range m {
    result = append(result, k)
}
	return result

}

// History of an item
func (bm *Blockmanager) GetItemHistory(itemId string, bcServer []Block) []Transaction {
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
