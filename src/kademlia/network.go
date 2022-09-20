package kademlia

import (
	"fmt"
	"log"
	"net"
	"time"
)

type Network struct {
	Kademlia     *Kademlia
	CurrentNode  *Contact
	RoutingTable *RoutingTable
}

func NewNetwork(node *Contact) *Network {
	return &Network{NewKademliaStruct(), node, NewRoutingTable(*node)}
}

func Listen(ip string, port int) {
	// TODO
}

// https://neo-ngd.github.io/NEO-Tutorial/en/5-network/2-Developing_a_NEO_ping_using_Golang.html
// TODO: FIX THIS!
func (network *Network) SendPingMessage(contact *Contact) bool {
	conn, err3 := net.Dial("udp4", contact.Address)
	if err3 != nil {
		log.Println(err3)
	}
	defer conn.Close()
	conn.Write([]byte("Test message123"))
	buffer := make([]byte, 4096)
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Println("\tResponse from server:", string(buffer[:n]))
	return true
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}
