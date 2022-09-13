package main

import (
	"github.com/PontusNorrby/D7024E-Kademlia"
)

func main() {
	//rt := NewRoutingTable(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000"))
	rt := d7024e.NewRoutingTable(d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000"))
}
