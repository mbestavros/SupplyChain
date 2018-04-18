# Blockmanager
Package that generates, verifies, and mines blocks
### Use cases:
  - Instantiate: <br>
    `	bm := Blockmanager{}`
  - Start a genesis block: <br>
    `	genBlock := bm.genesis()`
  - Generating a block: 
    - Generating a block now has mining which takes some undefined time (based on a difficulty constant) before actually adding a transaction to the blockchain <br>
    `myBlock := bm.generateBlock(oldBlock, <transaction struct>)`
  - Verifying a block is the successor to the old block: <br>
    `	ok := bm.isBlockValid(newBlock, oldBlock)`
