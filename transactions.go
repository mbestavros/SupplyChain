package transactions
//package main
//import "fmt"
type Action string


//Enum for all types of transactions
const (
	Create Action = "Create"
	Exchange Action = "Exchange"
	Consume Action = "Consume"
	Make Action = "Make"
	Split Action = "Split"
)

//Struct to represent all Transactions
type Transaction struct{
	Type Action
	OriginUser int64
	DestinationUser int64
	InitialTimestamp int
	FinalTimestamp int
	//Hash probably goes here as well @Sean

}

//func main(){
//	fmt.Println("dude nice")
//}
