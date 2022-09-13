package main

import (
	"fmt"
	"github.com/PontusNorrby/D7024E-Kademlia"
)

func main() {
	//rt := NewRoutingTable(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000"))
	rt := d7024e.NewRoutingTable(d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000"))

	rt.AddContact(d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8001"))
	rt.AddContact(d7024e.NewContact(d7024e.NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(d7024e.NewContact(d7024e.NewKademliaID("1111111200000000000000000000000000000000"), "localhost:8003"))
	rt.AddContact(d7024e.NewContact(d7024e.NewKademliaID("1111111300000000000000000000000000000000"), "localhost:8004"))
	rt.AddContact(d7024e.NewContact(d7024e.NewKademliaID("1111111400000000000000000000000000000000"), "localhost:8005"))
	rt.AddContact(d7024e.NewContact(d7024e.NewKademliaID("2111111400000000000000000000000000000000"), "localhost:8006"))

	contacts := rt.FindClosestContacts(d7024e.NewKademliaID("2111111400000000000000000000000000000000"), 20)
	for i := range contacts {
		fmt.Println(contacts[i].String())
	}

}
