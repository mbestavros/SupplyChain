To use the grouper, import the "grouper" package (this folder), and make a 
grouper object like this:

gr := Grouper{}

To start a network, use the StartNetwork function of Grouper as such:

// specify IP address, port, and your name
gr.StartNetwork("localhost", "8080", "ken")

To join an existing network, use:

// specify your friend's IP and port, then your own, then your name
gr.JoinNetwork("localhost", "8080", "localhost", "8081", "john")

Note that the StartNetwork and JoinNetwork commands are non-blocking, and start
an HTTP server in the background.

Finally, when you're exiting the program (you can capture a SIGINT event), 
make sure to call gr.Shutdown(), which will announce to the other nodes that 
this node is leaving the network and they can remove it from their list of
peers.

