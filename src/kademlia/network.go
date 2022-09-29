package kademlia

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

type Network struct {
	CurrentNode  *Contact
	RoutingTable *RoutingTable
}

func NewNetwork(node *Contact) *Network {
	return &Network{node, NewRoutingTable(*node)}
}

func (network *Network) Listen(ip string, port int) {
	addrToString := strings.Join([]string{ip, strconv.Itoa(port)}, ":")
	resolveUdpAddress, resolveUdpError := net.ResolveUDPAddr("udp4", addrToString)
	if resolveUdpError != nil {
		fmt.Println(resolveUdpError)
		panic(resolveUdpError)
	}
	listenUdpResponse, listenUdpError := net.ListenUDP("udp4", resolveUdpAddress)
	if listenUdpError != nil {
		//TODO: We end up here on everything other than first node...
		fmt.Println("error is", listenUdpError)
		return
	}

	fmt.Println("UDP server up and listening on", addrToString)

	defer func(listenUdpResponse *net.UDPConn) {
		closeError := listenUdpResponse.Close()
		if closeError != nil {
			return
		}
	}(listenUdpResponse)

	for {
		// wait for UDP client to connect
		buffer := make([]byte, 1024)
		n, readFromUDPAddress, readFromUdpError := listenUdpResponse.ReadFromUDP(buffer)

		if readFromUdpError != nil {
			log.Fatal(readFromUdpError)
		}

		fmt.Println("\tReceived from UDP client :", string(buffer[:n]))
		fmt.Println("Buffer: ", buffer[:n])
		message := getResponseMessage(buffer[:n], network)

		_, writeToUDPError := listenUdpResponse.WriteToUDP(message, readFromUDPAddress)

		if writeToUDPError != nil {
			log.Fatal(writeToUDPError)
		}
	}
}

// TODO: COMPLETE THIS
func getResponseMessage(message []byte, network *Network) []byte {
	messageList := strings.Split(string(message), " ")
	if messageList[0] == "Ping" {
		fmt.Println("Received Ping")
		body, err := json.Marshal(network.CurrentNode)
		if err != nil {
			log.Println(err)
			panic(err)
		}
		ex := extractContact([]byte(messageList[1]), network)
		if ex != nil {
			return ex
		}
		return body

	} else if messageList[0] == "FindContact" {
		//TODO

	} else if messageList[0] == "FindData" {
		//TODO

	} else if messageList[0] == "StoreMessage" {
		fmt.Println("Received StoreMessage")
		//Save the byte in the kademlia map... i guess?

	}
	return []byte("Error: Invalid RPC protocol")
}

// TODO: Which is better if statements or switch case?
//func getResponseMessage(message []byte, network *Network) []byte {
//	messageList := strings.Split(string(message), " ")
//	switch{
//	case messageList[0] == "Ping":
//		body, err := json.Marshal(network.CurrentNode)
//		if err != nil {
//			log.Println(err)
//			panic(err)
//		}
//		ex := extractContact([]byte(messageList[1]), network)
//		if ex != nil {
//			return ex
//		}
//		return body
//	}
//	return []byte("Error: Invalid RPC protocol")
//}

func extractContact(message []byte, network *Network) []byte {
	var contact *Contact
	err := json.Unmarshal(message, &contact)
	if err != nil {
		return nil
	}
	network.RoutingTable.AddContact(*contact)
	return nil
}

// https://neo-ngd.github.io/NEO-Tutorial/en/5-network/2-Developing_a_NEO_ping_using_Golang.html
// TODO: FIX THIS!
func (network *Network) SendPingMessage(contact *Contact) bool {
	conn, err3 := net.Dial("udp4", contact.Address)
	if err3 != nil {
		fmt.Println("test123")
		log.Println(err3)
	}
	defer conn.Close()

	// Message builder
	startMessage := []byte("Ping" + " ")
	body, err := json.Marshal(network.CurrentNode)
	if err != nil {
		log.Println(err)
		//panic(err)
	}
	message := append(startMessage, body...)
	conn.Write(message)

	buffer := make([]byte, 4096)
	//conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("test5")
		fmt.Println(err)
		return false
	}
	fmt.Println("\tResponse from server:", string(buffer[:n]))
	handlePingResponse(buffer[:n], network)
	return true
}

func handlePingResponse(message []byte, network *Network) {
	if string(message[:5]) == "Error" {
		log.Println(string(message))
		return
	} else {
		var contact *Contact
		json.Unmarshal(message, &contact)
		network.RoutingTable.AddContact(*contact)
	}
	// fmt.Println("ping response: ", network.routingTable)
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

/*func (network *Network) SendStoreMessage(data []byte) bool {
	// TODO
	//Maybe a check if the data is too big to store?
	conn, err3 := net.Dial("udp4", network.CurrentNode.Address)
	if err3 != nil {
		log.Println(err3)
	}
	defer conn.Close()
	// Message builder
	startMessage := []byte("StoreMessage" + " ")
	body := data
	message := append(startMessage, body...)
	conn.Write(message)

	buffer := make([]byte, 4096)
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Println("\tResponse from server:", string(buffer[:n]))
	//Ska implementeras en hantering av store efter när algorithmen för att hitta
	//vart value ska sparas så rätt bucket uppdateras.
	//handleStoreResponse(buffer[:n], network)
	return true
}

func handleStoreResponse(message []byte, network *Network) {
	if string(message[:5]) == "Error" {
		log.Println(string(message))
		return
	} else {
		var contact *Contact
		json.Unmarshal(message, &contact)
		network.RoutingTable.AddContact(*contact)
	}
}*/
