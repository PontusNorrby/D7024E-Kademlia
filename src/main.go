package main

import (
	"fmt"
	"github.com/PontusNorrby/D7024E-Kademlia/src/kademlia"
	"log"
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

	fmt.Println("My IP: ", localIp)
	fmt.Println("Base IP: ", BaseIp)

	newNetwork := kademlia.NewNetwork(&selfContact)

	if GetOutboundIP().String() != BaseIp {
		selfContact.Address = localIp.String()
		print(selfContact.Address + "\n")

		newNetwork.RoutingTable.AddContact(baseContact)
		print(newNetwork.RoutingTable)
	}

	if localIp.String() != BaseIp {
		if !newNetwork.SendPingMessage(&baseContact) {
			//If you end up here the baseNode is dead.
			panic("Can't connect to the network base node is down")
		}
	}

	newNetwork.Listen(BaseIp, Port)
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
