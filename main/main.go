package main

import (
	"fmt"
	d7024e "github.com/PontusNorrby/D7024E-Kademlia"
	"log"
	"net"
)

var (
	BaseIp string = "172.20.128.2"
)

func main() {
	id := d7024e.NewRandomKademliaID()
	LocalIp := GetOutboundIP()
	selfContact := d7024e.NewContact(id, LocalIp.String())
	newRouteTable := d7024e.NewRoutingTable(selfContact)
	if selfContact.String() == BaseIp {
		//First node in the network
	}
	//A new node in the network
	//func join network..
	fmt.Printf("", newRouteTable)
}

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
