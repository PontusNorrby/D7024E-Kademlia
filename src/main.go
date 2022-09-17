package main

import (
	"fmt"
	"github.com/PontusNorrby/D7024E-Kademlia/src/kademlia"
	"log"
	"net"
)

var (
	BaseIp string = "172.20.128.2"
)

func main() {
	id := kademlia.NewRandomKademliaID()
	LocalIp := GetOutboundIP()
	selfContact := kademlia.NewContact(id, LocalIp.String())
	newRouteTable := kademlia.NewRoutingTable(selfContact)
	if selfContact.String() == BaseIp {
		//First node in the network
	}
	//A new node in the network
	//func join network.
	fmt.Printf("", newRouteTable)
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
