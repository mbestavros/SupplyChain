# SupplyChain
Final project for CS651: Distributed Systems. Blockchain for supply-chain verification.

Setup:
* Set your GOPATH variable to be the root folder of this repo, and go to that folder. You can do this by navigating to the main directory and running:  
`GOPATH=\$(pwd)`

* Run the following to get required packages:  
`go get github.com/rs/xid`  
`go get github.com/sirupsen/logrus`


Build, Run, and Test:  
* Whole project  
`go build cli`  
`./cli`  
* Individual packages  
`cd <package-name>` 
To run:  
`go build <package-name>`  
`./<package-name>`  
To test:  
`go test`

## Debugging

We're using the logrus library. To set the log level to Debug, for example, write `log.SetLevel(log.DebugLevel)`. To then log something to Debug, use `log.Debug("Debug statement!")`.
