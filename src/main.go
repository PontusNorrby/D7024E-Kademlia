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
	//id := kademlia.NewRandomKademliaID()
	//LocalIp := GetOutboundIP()
	//selfContact := kademlia.NewContact(id, LocalIp.String())
	//newRouteTable := kademlia.NewRoutingTable(selfContact)

	//This is the node every node is going to join
	target := kademlia.NewRandomKademliaID()
	contact := kademlia.NewContact(target, BaseIp+":"+strconv.Itoa(Port))

	LocalIp := GetOutboundIP()
	// The current node, aka this node
	currentContact := kademlia.NewContact(kademlia.NewRandomKademliaID(), LocalIp.String())
	network := kademlia.NewNetwork(&currentContact)

	fmt.Println("My IP: ", LocalIp)

	ping := network.SendPingMessage(&contact)
	fmt.Println(ping)
	fmt.Println("contact.String(): ", contact.String())
	fmt.Println("BaseIp: ", BaseIp)
	if contact.String() == BaseIp {
		fmt.Println("Node IP: ", LocalIp.String())
	}

	fmt.Println(network.RoutingTable)

	//A new node in the network
	//func join network.
	fmt.Printf("", network.RoutingTable)
}

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
