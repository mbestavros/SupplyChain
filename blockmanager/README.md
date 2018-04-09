# Blockmanager
Package that generates and verifies blocks
### Use cases:
  - Instantiate: <br>
    `	bm := Blockmanager{}`
  - Generating a block: <br>
    `myBlock := bm.generateBlock(oldBlock, <transaction string>)`
  - Verifying a block is the successor to the old block: <br>
    `	ok := bm.isBlockValid(newBlock, oldBlock)`

### Note
  - In future iterations, we should probably change the Transaction property of Block to a struct or something else, currently it's just a string
