# CS651-Spr18-Final-Project
Final project for CS651: Distributed Systems. Blockchain mixed with distributed systems.

Setting it up:

* Set your GOPATH variable to be the root folder of this repo, go to that folder
* run "go build pkg-name" where pkg-name is the package name where your main func is (for example: "go build cli" to run our main function)
* "./pkg-name" will execute the main func (for us, that's "./cli")

A few notes for us:

* When you're working on a part of the code, create a folder in "src" for it. That's your package. We have packages: grouper, blockmanager, server, and cli).
* Make a file pkg-name.go in that folder, the first line of that code should be "package pkg-name"

To run the program, type "go build cli" and then run "./cli" to launch the command line interface.

To test your code, create a TestFunc in it, then run "go test grouper -v" (if you're testing grouper). Look in the grouper package for a description of how to test individual packages.


