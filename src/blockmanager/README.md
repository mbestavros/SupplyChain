# Blockmanager
Package that generates, verifies, and mines blocks
### Use cases:
  - Instantiate: <br>
    `	bm := Blockmanager{}`
  - Start a genesis block: <br>
    `	genBlock := bm.Genesis()`
  - Generating a block: 
    - Generating a block now has mining which takes some undefined time (based on a difficulty constant) before actually adding a transaction to the blockchain <br>
    `myBlock := bm.GenerateBlock(oldBlock, <TransactionProvider>)`
  - Verifying a block is the successor to the old block: <br>
    `	ok := bm.isBlockValid(newBlock, oldBlock)`

# Transactions
Golang does not have traditional inheritance like Java. There's some hacks to get polymorphic behavior.

At the top of the hierarchy, there is a `TransactionProvider` interface. Then, we have a "base" class `Transaction` struct that contains `TimeTransacted` and `TransactionType` for all subtransaction types to inherit. The "child" classes (i.e. `CreateTransaction`, etc) inherit the `Transaction` struct's attributes and can add their own attributes as well.

The base `Transaction` struct implements the `TransactionProvider` interface (`GetTransaction()`), which means that all subtransaction types (i.e. Create, Exchange, etc) are "subclasses" of type `TransactionProvider` by embedding the `Transaction` struct within their fields.

Then, the `GenerateBlock` takes a generic `TransactionProvider` as a parameter, allowing it to accept any type of Transaction

#### Resources
- https://stackoverflow.com/questions/26027350/go-interface-fields
- https://medium.com/golangspec/interfaces-in-go-part-iii-61f5e7c52fb5
