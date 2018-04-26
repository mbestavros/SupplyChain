package blockmanager

type Action string

//Enum for all types of transactions
const (
	Create   Action = "Create"
	Exchange Action = "Exchange"
	Consume  Action = "Consume"
	Make     Action = "Make"
	Split    Action = "Split"
)

// define Transaction as an interface to allow for polymorphism
type TransactionProvider interface {
	GetTransaction() *Transaction
}

// base Transaction class
type Transaction struct {
	TransactionType Action
	TimeTransacted  int64
}

// Transaction struct "implements" TransactionProvider interface
func (t *Transaction) GetTransaction() *Transaction {
	return t
}

type CreateTransaction struct {
	// create transaction origin user will be set to "" because a material is being created
	Transaction       // embedded struct
	OriginUserId      string
	DestinationUserId string
	ItemId            string
	ItemName          string
}

type ExchangeTransaction struct {
	// exchange transaction is a canonical transaction, give an item from one origin to new destination
	Transaction       // embedded interface
	OriginUserId      string
	DestinationUserId string
	ItemId            string
	ItemName          string
}

type ConsumeTransaction struct {
	// even though originUser == destinationUser, keeping destination attribute for history function / provenance tracking later
	Transaction       // embedded interface
	OriginUserId      string
	DestinationUserId string
	ItemId            string
	ItemName          string
}

type MakeTransaction struct {
	// many input items -> one output item. originUser == destinationUser
	Transaction       // embedded interface
	OriginUserId      string
	DestinationUserId string
	InputItemIds      []string
	InputItemNames    []string
	OutputItemId      string
	OutputItemName    string
}

type SplitTransaction struct {
	// 1 origin user -> many destination users. If user is splitting item with themselves, then DestinationUserIds can just be a list of size one
	Transaction        // embedded interface
	OriginUserId       string
	DestinationUserIds []string
	InputItemId        string
	InputItemName      string
	OutputItemIds      []string
	OutputItemNames    []string
}
