package main

import (
	"github.com/PontusNorrby/D7024E-Kademlia/src/kademlia"
	"log"
	"math/rand"
	"net"
	"strconv"
)

var (
	BaseIp string = "172.20.0.2"
	Port   int    = 3000
)

func main() {
	//The newly created node gets a random id
	selfId := kademlia.NewRandomKademliaID()
	localIp := GetOutboundIP()

	//Creates contacts for both the new node and the base node
	selfContact := kademlia.NewContact(selfId, "")
	baseContact := kademlia.NewContact(kademlia.NewRandomKademliaID(), BaseIp+":"+strconv.Itoa(Port))

	newNetwork := kademlia.NewNetwork(&selfContact)

	LocalIp := GetOutboundIP()
	// The current node, aka this node
	currentContact := kademlia.NewContact(kademlia.NewRandomKademliaID(), LocalIp.String())
	network := kademlia.NewNetwork(&currentContact)

	rand.Seed()

	if !newNetwork.SendPingMessage(&baseContact) {
		//If you end up here the baseNode is dead?
		panic("Can't connect to the network")
	}
	//Om pingen returnerar korrekt slumpad id så kan vi uppdatera egna bucketen.
	//update own bucket since we got the correct random id back.

	//Här görs allt som ska hända när det är en ny nod och inte basnoden
	if GetOutboundIP().String() != BaseIp {
		selfContact.Address = localIp.String()
		print(selfContact.Address + "\n")
		newNetwork.RoutingTable.AddContact(baseContact)
		//newNetwork.RoutingTable.FindClosestContacts(selfContact.ID, 20)
		print(newNetwork.RoutingTable)
	}
	//Funktionalitet för basnoden

}

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			print("Closed connection")
		}
	}(conn)

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
