package main

import (
	"fmt"
	"github.com/PontusNorrby/D7024E-Kademlia/src/cli"
	"github.com/PontusNorrby/D7024E-Kademlia/src/kademlia"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"
)

var (
	BaseIp         = "172.20.0.2"
	Port           = 3000
	kademliaStruct *kademlia.Kademlia
	baseContact    kademlia.Contact
)

func main() {
	rand.Seed(time.Now().UnixNano())

	//The newly created node gets a random id
	selfId := kademlia.NewRandomKademliaID()
	localIP := GetOutboundIP()

	//Creates contacts for the new node
	selfContact := kademlia.NewContact(selfId, "")

	if localIP.String() == BaseIp {
		baseContact = kademlia.NewContact(kademlia.NewRandomKademliaID(), BaseIp+":"+strconv.Itoa(Port))
		kademliaPing := kademlia.NewNetwork(&selfContact).SendPingMessage(&baseContact)
		if kademliaPing {
			Port = rand.Intn(65535-1024) + 1024
		}
		selfContact = kademlia.NewContact(kademlia.NewRandomKademliaID(), localIP.String()+":"+strconv.Itoa(Port))
		kademliaStruct = kademlia.NewKademliaStruct(kademlia.NewNetwork(&baseContact))

	} else {
		baseContact = kademlia.NewContact(kademlia.NewRandomKademliaID(), BaseIp+":"+strconv.Itoa(Port))
		Port = rand.Intn(65535-1024) + 1024
		selfContact = kademlia.NewContact(kademlia.NewRandomKademliaID(), localIP.String()+":"+strconv.Itoa(Port))
		kademliaStruct = kademlia.NewKademliaStruct(kademlia.NewNetwork(&baseContact))
		success := kademliaStruct.Network.SendPingMessage(&baseContact)
		if !success {
			panic("Can't connect to baseNode")
		}
	}

	go kademliaStruct.Network.Listen(localIP.String(), Port, kademliaStruct)
	cli.StartCLI(exitNode, kademliaStruct)
	time.Sleep(1 * time.Second)
	kademliaStruct.Network.SendFindContactMessage(&baseContact, selfContact.ID)
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

func exitNode() {
	fmt.Println("Shutting down the node...")
	os.Exit(0)
}
