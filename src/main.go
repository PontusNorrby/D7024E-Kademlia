package main

import (
	"github.com/PontusNorrby/D7024E-Kademlia/src/kademlia"
	"log"
	"math/rand"
	"net"
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
	//fmt.Println("Base IP: ", BaseIp)

	//testStoreValue := []byte("HelloWorld")
	/*newKademlia := kademlia.NewKademliaStruct(newNetwork)
	newKademlia.Store(testStoreValue)*/

	if localIP.String() == BaseIp {

		kademliaPing := kademlia.NewNetwork(&selfContact).SendPingMessage(&baseContact)
		if kademliaPing {
			Port = rand.Intn(65535-1024) + 1024
		}
		selfContact = kademlia.NewContact(kademlia.NewRandomKademliaID(), localIP.String()+":"+strconv.Itoa(Port))
		baseContact = kademlia.NewContact(kademlia.NewRandomKademliaID(), localIP.String()+":"+strconv.Itoa(Port))
		kademliaStruct = kademlia.NewKademliaStruct(kademlia.NewNetwork(&selfContact))
	} else {
		//baseContact = kademlia.NewContact(kademlia.NewRandomKademliaID(), localIP.String()+":"+strconv.Itoa(Port))
		selfContact = kademlia.NewContact(kademlia.NewRandomKademliaID(), localIP.String()+":"+strconv.Itoa(Port))
		kademliaStruct = kademlia.NewKademliaStruct(kademlia.NewNetwork(&selfContact))
		if !kademliaStruct.Network.SendPingMessage(&baseContact) {
			panic("Can't connect to baseNode")
		}
		Port = rand.Intn(65535-1024) + 1024
	}

	//newNetwork.SendStoreMessage(testStoreValue)

	go kademliaStruct.Network.Listen(localIP.String(), Port, kademliaStruct)

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
