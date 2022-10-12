package kademlia

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

type Network struct {
	CurrentNode  *Contact
	RoutingTable *RoutingTable
}

func NewNetwork(node *Contact) *Network {
	return &Network{node, NewRoutingTable(*node)}
}

func (network *Network) Listen(ip string, port int, kademliaStruct *Kademlia) {
	addrToString := strings.Join([]string{ip, strconv.Itoa(port)}, ":")
	resolveUdpAddress, resolveUdpError := net.ResolveUDPAddr("udp4", addrToString)
	if resolveUdpError != nil {
		//fmt.Println(resolveUdpError)
		panic(resolveUdpError)
	}
	listenUdpResponse, listenUdpError := net.ListenUDP("udp4", resolveUdpAddress)
	if listenUdpError != nil {
		//TODO: We end up here on everything other than first node...
		//fmt.Println("error is", listenUdpError)
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
		buffer := make([]byte, 4096)
		n, readFromUDPAddress, readFromUdpError := listenUdpResponse.ReadFromUDP(buffer)

		if readFromUdpError != nil {
			log.Fatal(readFromUdpError)
		}

		//fmt.Println("\tReceived from UDP client :", string(buffer[:n]))
		//fmt.Println("Buffer: ", buffer[:n])
		message := getResponseMessage(buffer[:n], network, kademliaStruct)

		_, writeToUDPError := listenUdpResponse.WriteToUDP(message, readFromUDPAddress)

		if writeToUDPError != nil {
			log.Fatal(writeToUDPError)
		}
	}
}

// TODO: COMPLETE THIS
func getResponseMessage(message []byte, network *Network, kademliaStruct *Kademlia) []byte {
	messageList := strings.Split(string(message), " ")
	if messageList[0] == "Ping" {
		//fmt.Println("Received Ping")
		body, err := json.Marshal(network.CurrentNode)
		if err != nil {
			//log.Println(err)
			panic(err)
		}
		extraction := extractContact([]byte(messageList[1]), network)
		if extraction != nil {
			return extraction
		}
		return body

	} else if messageList[0] == "FindContact" {
		var id *KademliaID
		UnmarshalError := json.Unmarshal([]byte(messageList[1]), &id)
		if UnmarshalError != nil {
			println("Error is ", UnmarshalError)
			return nil
		}
		extraction := extractContact([]byte(messageList[2]), network)
		if extraction != nil {
			fmt.Println(extraction)
			return extraction
		}
		closestNodes := network.RoutingTable.FindClosestContacts(id, 20)
		closestNodes = append(closestNodes, *network.CurrentNode)
		body, _ := json.Marshal(closestNodes)
		return body

	} else if messageList[0] == "FindData" {
		var hash *KademliaID
		json.Unmarshal([]byte(messageList[1]), &hash)
		ex := extractContact([]byte(messageList[2]), network)
		if ex != nil {
			fmt.Println(ex)
			return ex
		}
		lookupValue := kademliaStruct.LookupData(*hash)
		if lookupValue != nil {
			body, _ := json.Marshal(network.CurrentNode)
			return []byte("VALUE;" + string(lookupValue) + " " + string(body))
		}
		resClosestNodes := network.RoutingTable.FindClosestContacts(hash, 4096)
		resClosestNodes = append(resClosestNodes, *network.CurrentNode)
		body, _ := json.Marshal(resClosestNodes)
		return []byte("CONT" + string(body))

	} else if messageList[0] == "StoreMessage" {
		//fmt.Println("Received StoreMessage")
		var storeData *[]byte
		json.Unmarshal([]byte(messageList[1]), &storeData)
		kademliaStruct.store(*storeData)
		ex := extractContact([]byte(messageList[2]), network)
		if ex != nil {
			return ex
		}
		body, _ := json.Marshal(network.CurrentNode)
		return body
	}
	return []byte("Error: Invalid RPC protocol")
}

func extractContact(message []byte, network *Network) []byte {
	var contact *Contact
	err := json.Unmarshal(message, &contact)
	if err != nil {
		return nil
	}
	network.RoutingTable.AddContact(*contact)
	return nil
}

// SendPingMessage https://neo-ngd.github.io/NEO-Tutorial/en/5-network/2-Developing_a_NEO_ping_using_Golang.html
func (network *Network) SendPingMessage(contact *Contact) bool {
	//fmt.Println("123,")
	conn, err3 := net.Dial("udp4", contact.Address)
	if err3 != nil {
		log.Println(err3)
	}

	defer conn.Close()

	//Message builder
	startMessage := []byte("Ping" + " ")
	body, err5 := json.Marshal(network.CurrentNode)
	if err5 != nil {
		//fmt.Println("test5")
		log.Println(err5)
	}
	message := append(startMessage, body...)

	conn.Write(message)
	buffer := make([]byte, 4096)
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	n, err := conn.Read(buffer)
	if err != nil {
		return false
	}
	// fmt.Println("\tResponse from server:", string(buffer[:n]))
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

func (network *Network) SendFindContactMessage(contact *Contact, searchID *KademliaID) []Contact {
	conn, err3 := net.Dial("udp4", contact.Address)
	if err3 != nil {
		log.Println(err3)
	}
	defer conn.Close()

	//Message builder
	body, _ := json.Marshal(searchID)
	startMessage := []byte("FindContact" + " " + string(body))
	body, _ = json.Marshal(network.CurrentNode)

	message := append(startMessage, body...)

	conn.Write(message)
	buffer := make([]byte, 4096)
	//fmt.Println("test")
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	n, err := conn.Read(buffer)
	if err != nil {
		return nil
	}
	//fmt.Println("\tResponse from server:", string(buffer[:n]))
	return handleFindContactResponse(buffer[:n], network)
}

func (network *Network) SendFindDataMessage(hash string, contact *Contact) string {
	conn, err3 := net.Dial("udp4", contact.Address)
	if err3 != nil {
		log.Println(err3)
		return "ERROR"
	}
	defer conn.Close()

	//Message builder
	body, _ := json.Marshal(hash)
	startMessage := []byte("FindData" + " " + string(body) + " ")
	body2, _ := json.Marshal(network.CurrentNode)
	message := append(startMessage, body2...)

	conn.Write(message)
	buffer := make([]byte, 4096)
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	n, err := conn.Read(buffer)
	if err != nil {
		return "ERROR"
	}
	// fmt.Println("\tResponse from server:", string(buffer[:n]))
	return handleSendDataResponse(buffer[:n], network)

}

func handleSendDataResponse(message []byte, network *Network) string {
	if string(message[:5]) == "Error" {
		log.Println(string(message))
		return string(message)
	} else {
		if string(message[:4]) == "VALU" {
			resp := strings.Split(string(message[5:]), " ")
			var contact *Contact
			json.Unmarshal([]byte(resp[1]), &contact)
			if contactUsability(contact, network) {
				network.RoutingTable.AddContact(*contact)
			}
			return resp[0]
		}
		var foundContact []Contact
		json.Unmarshal(message, &foundContact)
		for _, foundContact := range foundContact {
			if contactUsability(&foundContact, network) {
				network.SendPingMessage(&foundContact)
			}
		}
		return ""
	}
}

func handleFindContactResponse(message []byte, network *Network) []Contact {
	if string(message[:5]) == "Error" {
		log.Println(string(message))
		return nil
	} else {
		var foundContacts []Contact
		var usableContacts []Contact
		json.Unmarshal(message, &foundContacts)
		for _, foundContact := range foundContacts {
			if contactUsability(&foundContact, network) {
				if network.SendPingMessage(&foundContact) {
					usableContacts = append(usableContacts, foundContact)
				}
			}
		}
		return usableContacts
	}
}

func contactUsability(contact *Contact, network *Network) bool {
	//Returns false if address empty or checking on it's own ID.
	if contact.Address == "" || contact.ID == network.CurrentNode.ID {
		return false
	}
	return true
}

func (network *Network) SendStoreMessage(data []byte, contact *Contact, kademlia *Kademlia) bool {
	conn, err3 := net.Dial("udp4", contact.Address)
	if err3 != nil {
		log.Println(err3)
		return false
	}
	defer conn.Close()
	//Message builder
	body, _ := json.Marshal(data)
	message := []byte("StoreMessage" + " " + string(body))
	body, _ = json.Marshal(network.CurrentNode)
	message = append(message, body...)

	conn.Write(message)
	buffer := make([]byte, 4096)
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	n, err := conn.Read(buffer)
	if err != nil {
		return false
	}
	// fmt.Println("\tResponse from server:", string(buffer[:n]))
	handleStoreResponse(buffer[:n], network)
	return true
}

func handleStoreResponse(message []byte, network *Network) {
	if string(message[:5]) == "Error" {
		log.Println(string(message))
		return
	} else {
		var storeContact *Contact
		json.Unmarshal(message, &storeContact)
		if contactUsability(storeContact, network) {
			network.RoutingTable.AddContact(*storeContact)
		}
	}
}
