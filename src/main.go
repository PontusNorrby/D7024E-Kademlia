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
	BaseIp         string = "172.20.0.2"
	Port           int    = 3000
	kademliaStruct *kademlia.Kademlia
	//OtherPort int    = 3001
)

func main() {
	rand.Seed(time.Now().UnixNano())

	//The newly created node gets a random id
	selfId := kademlia.NewRandomKademliaID()
	localIP := GetOutboundIP()

	//Creates contacts for both the new node and the base node
	selfContact := kademlia.NewContact(selfId, "")
	baseContact := kademlia.NewContact(kademlia.NewRandomKademliaID(), BaseIp+":"+strconv.Itoa(Port))

	if localIP.String() == BaseIp {
		kademliaPing := kademlia.NewNetwork(&selfContact).SendPingMessage(&baseContact)
		if kademliaPing {
			Port = rand.Intn(65535-1024) + 1024
		}
		selfContact = kademlia.NewContact(kademlia.NewRandomKademliaID(), localIP.String()+":"+strconv.Itoa(Port))
		baseContact = kademlia.NewContact(kademlia.NewRandomKademliaID(), localIP.String()+":"+strconv.Itoa(Port))
		kademliaStruct = kademlia.NewKademliaStruct(kademlia.NewNetwork(&selfContact))
	} else {
		selfContact = kademlia.NewContact(kademlia.NewRandomKademliaID(), localIP.String()+":"+strconv.Itoa(Port))
		kademliaStruct = kademlia.NewKademliaStruct(kademlia.NewNetwork(&selfContact))
		if !kademliaStruct.Network.SendPingMessage(&baseContact) {
			panic("Can't connect to baseNode")
		}
		Port = rand.Intn(65535-1024) + 1024
	}

	go kademliaStruct.Network.Listen(localIP.String(), Port, kademliaStruct)
	go cli.StartCLI(nodeShutDown, kademliaStruct)
	//b := []byte("ABC")
	//kademliaStruct.StoreValue(b)

	for {
		time.Sleep(5 * time.Second)
	}

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

func nodeShutDown() {
	fmt.Println("Shutting down the node...")
	os.Exit(0)
}
