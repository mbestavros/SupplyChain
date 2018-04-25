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
type Transaction interface {
}

type CreateTransaction struct {
	// create transaction has no origin user because a material is being created
	Transaction       // embedded interface
	TransactionType   Action
	DestinationUserId int64
	ItemId            string
	ItemName          string
	TimeTransacted    int64
}

type ExchangeTransaction struct {
	// exchange transaction is a canonical transaction, give an item from one origin to new destination
	Transaction       // embedded interface
	TransactionType   Action
	OriginUserId      int64
	DestinationUserId int64
	ItemId            string
	ItemName          string
	TimeTransacted    int64
}

type ConsumeTransaction struct {
	// even though originUser == destinationUser, keeping destination attribute for history function / provenance tracking later
	Transaction       // embedded interface
	TransactionType   Action
	OriginUserId      int64
	DestinationUserId int64
	ItemId            string
	ItemName          string
	TimeTransacted    int64
}

type MakeTransaction struct {
	// many input items -> one output item. originUser == destinationUser
	Transaction       // embedded interface
	TransactionType   Action
	OriginUserId      int64
	DestinationUserId int64
	InputItemIds      []string
	InputItemNames    []string
	OutputItemId      string
	OutputItemName    string
	TimeTransacted    int64
}

type SplitTransaction struct {
	// 1 origin user -> many destination users. If user is splitting item with themselves, then DestinationUserIds can just be a list of size one
	Transaction        // embedded interface
	TransactionType    Action
	OriginUserId       int64
	DestinationUserIds []int64
	InputItemId        string
	InputItemName      string
	OutputItemIds      []string
	OutputItemNames    []string
	TimeTransacted     int64
}
