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

// base Transaction class
type Transaction struct {
	TransactionType Action
	TimeTransacted  int64

	Cr CreateTransaction
	Ex ExchangeTransaction
	Co ConsumeTransaction
	Ma MakeTransaction
	Sp SplitTransaction
}

type CreateTransaction struct {
	// create transaction origin user will be set to "" because a material is being created
	OriginUserId      string
	DestinationUserId string
	ItemId            string
	ItemName          string
}

type ExchangeTransaction struct {
	// exchange transaction is a canonical transaction, give an item from one origin to new destination
	OriginUserId      string
	DestinationUserId string
	ItemId            string
	ItemName          string
}

type ConsumeTransaction struct {
	// even though originUser == destinationUser, keeping destination attribute for history function / provenance tracking later
	OriginUserId      string
	DestinationUserId string
	ItemId            string
	ItemName          string
}

type MakeTransaction struct {
	// many input items -> one output item. originUser == destinationUser
	OriginUserId      string
	DestinationUserId string
	InputItemIds      []string
	InputItemNames    []string
	OutputItemId      string
	OutputItemName    string
}

type SplitTransaction struct {
	// 1 origin user -> many destination users. If user is splitting item with themselves, then DestinationUserIds can just be a list of size one
	OriginUserId       string
	DestinationUserIds []string
	InputItemId        string
	InputItemName      string
	OutputItemIds      []string
	OutputItemNames    []string
}
