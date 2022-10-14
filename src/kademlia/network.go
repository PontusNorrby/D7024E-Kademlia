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
		panic(resolveUdpError)
	}
	listenUdpResponse, listenUdpError := net.ListenUDP("udp4", resolveUdpAddress)
	if listenUdpError != nil {
		return
	}

	//fmt.Println("UDP server up and listening on", addrToString)

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
		message := getResponseMessage(buffer[:n], network, kademliaStruct)

		_, writeToUDPError := listenUdpResponse.WriteToUDP(message, readFromUDPAddress)

		if writeToUDPError != nil {
			log.Fatal(writeToUDPError)
		}
	}
}

func getResponseMessage(message []byte, network *Network, kademliaStruct *Kademlia) []byte {
	messageList := strings.Split(string(message), " ")
	if messageList[0] == "Ping" {
		marshalBody, marshalError := json.Marshal(network.CurrentNode)
		if marshalError != nil {
			panic(marshalError)
		}
		extraction := extractContact([]byte(messageList[1]), network)
		if extraction != nil {
			return extraction
		}
		return marshalBody

	} else if messageList[0] == "FindContact" {
		var id *KademliaID
		UnmarshalError := json.Unmarshal([]byte(messageList[1]), &id)
		if UnmarshalError != nil {
			println("Error is: ", UnmarshalError)
			return nil
		}
		extraction := extractContact([]byte(messageList[2]), network)
		if extraction != nil {
			fmt.Println(extraction)
			return extraction
		}
		closestNodes := network.RoutingTable.FindClosestContacts(id, 20)
		closestNodes = append(closestNodes, *network.CurrentNode)
		marshalBody, _ := json.Marshal(closestNodes)
		return marshalBody

	} else if messageList[0] == "FindData" {
		var hash *KademliaID
		unmarshalError := json.Unmarshal([]byte(messageList[1]), &hash)
		if unmarshalError != nil {
			return nil
		}
		extraction := extractContact([]byte(messageList[2]), network)
		if extraction != nil {
			fmt.Println(extraction)
			return extraction
		}
		lookupValue := kademliaStruct.LookupData(*hash)
		if lookupValue != nil {
			marshalBody, _ := json.Marshal(network.CurrentNode)
			return []byte("VALUE" + string(lookupValue) + " " + string(marshalBody))
		}
		resClosestNodes := network.RoutingTable.FindClosestContacts(hash, 20)
		resClosestNodes = append(resClosestNodes, *network.CurrentNode)
		marshalBody, _ := json.Marshal(resClosestNodes)
		return []byte("CONT" + string(marshalBody))

	} else if messageList[0] == "StoreMessage" {
		var storeData *[]byte
		unmarshalError := json.Unmarshal([]byte(messageList[1]), &storeData)
		if unmarshalError != nil {
			return nil
		}
		kademliaStruct.storeDataHelp(*storeData)
		extraction := extractContact([]byte(messageList[2]), network)
		if extraction != nil {
			return extraction
		}
		marshalBody, _ := json.Marshal(network.CurrentNode)
		return marshalBody
	}
	return []byte("Error: Invalid RPC protocol")
}

func extractContact(message []byte, network *Network) []byte {
	var contact *Contact
	unmarshalError := json.Unmarshal(message, &contact)
	if unmarshalError != nil {
		return nil
	}
	if !contactUsability(contact, network) {
		return []byte("Error: Invalid contact information")
	}
	network.RoutingTable.AddContact(*contact)
	return nil
}

// SendPingMessage https://neo-ngd.github.io/NEO-Tutorial/en/5-network/2-Developing_a_NEO_ping_using_Golang.html
func (network *Network) SendPingMessage(contact *Contact) bool {
	netDialConnection, netDialError := net.Dial("udp4", contact.Address)
	if netDialError != nil {
		log.Println(netDialError)
		return false
	}
	defer func(netDialConnection net.Conn) {
		closeError := netDialConnection.Close()
		if closeError != nil {

		}
	}(netDialConnection)
	startMessage := []byte("Ping" + " ")
	marshalBody, _ := json.Marshal(network.CurrentNode)
	message := append(startMessage, marshalBody...)
	_, writeError := netDialConnection.Write(message)
	if writeError != nil {
		return false
	}
	buffer := make([]byte, 4096)
	setReadDeadlineError := netDialConnection.SetReadDeadline(time.Now().Add(2 * time.Second))
	if setReadDeadlineError != nil {
		return false
	}
	end, readError := netDialConnection.Read(buffer)
	if readError != nil {
		return false
	}
	handlePingResponse(buffer[:end], network)
	return true
}

func handlePingResponse(message []byte, network *Network) {
	if string(message[:5]) == "Error" {
		log.Println(string(message))
		return
	} else {
		var contact *Contact
		unmarshalError := json.Unmarshal(message, &contact)
		if unmarshalError != nil {
			return
		}
		if contactUsability(contact, network) {
			network.RoutingTable.AddContact(*contact)
		}
	}
}

func (network *Network) SendFindContactMessage(contact *Contact, searchID *KademliaID) []Contact {
	netDialConnection, netDialError := net.Dial("udp4", contact.Address)
	if netDialError != nil {
		log.Println(netDialError)
		return nil
	}
	defer func(netDialConnection net.Conn) {
		closeError := netDialConnection.Close()
		if closeError != nil {

		}
	}(netDialConnection)
	marshalBody, _ := json.Marshal(searchID)
	startMessage := []byte("FindContact" + " " + string(marshalBody) + " ")
	newMarshalBody, _ := json.Marshal(network.CurrentNode)
	message := append(startMessage, newMarshalBody...)
	_, writeError := netDialConnection.Write(message)
	if writeError != nil {
		return nil
	}
	buffer := make([]byte, 4096)
	setReadDeadlineError := netDialConnection.SetReadDeadline(time.Now().Add(2 * time.Second))
	if setReadDeadlineError != nil {
		return nil
	}
	end, readError := netDialConnection.Read(buffer)
	if readError != nil {
		return nil
	}
	return handleFindContactResponse(buffer[:end], network)
}

func handleFindContactResponse(message []byte, network *Network) []Contact {
	if string(message[:5]) == "Error" {
		log.Println(string(message))
		return nil
	} else {
		var foundContacts []Contact
		var usableContacts []Contact
		unmarshalError := json.Unmarshal(message, &foundContacts)
		if unmarshalError != nil {
			return nil
		}
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

func (network *Network) SendFindDataMessage(value *KademliaID, contact *Contact) string {
	netDialConnection, netDialError := net.Dial("udp4", contact.Address)
	if netDialError != nil {
		log.Println(netDialError)
		return "ERROR"
	}
	defer func(netDialConnection net.Conn) {
		closeError := netDialConnection.Close()
		if closeError != nil {

		}
	}(netDialConnection)
	marshalBody, _ := json.Marshal(value)
	startMessage := []byte("FindData" + " " + string(marshalBody) + " ")
	newMarshalBody, _ := json.Marshal(network.CurrentNode)
	message := append(startMessage, newMarshalBody...)
	_, writeError := netDialConnection.Write(message)
	if writeError != nil {
		return ""
	}
	buffer := make([]byte, 4096)
	SetReadDeadlineError := netDialConnection.SetReadDeadline(time.Now().Add(2 * time.Second))
	if SetReadDeadlineError != nil {
		return ""
	}
	end, readError := netDialConnection.Read(buffer)
	if readError != nil {
		return "ERROR"
	}
	return handleSendDataResponse(buffer[:end], network)

}

func handleSendDataResponse(message []byte, network *Network) string {
	if string(message[:5]) == "Error" {
		log.Println(string(message))
		return string(message)
	} else {
		if string(message[:5]) == "VALUE" {
			resp := strings.Split(string(message[5:]), " ")
			var contact *Contact
			unmarshalError := json.Unmarshal([]byte(resp[1]), &contact)
			if unmarshalError != nil {
				return ""
			}
			if contactUsability(contact, network) {
				network.RoutingTable.AddContact(*contact)
			}
			return resp[0]
		}
		var foundContact []Contact
		unmarshalError := json.Unmarshal(message, &foundContact)
		if unmarshalError != nil {
			return ""
		}
		for _, foundContact := range foundContact {
			if contactUsability(&foundContact, network) {
				network.SendPingMessage(&foundContact)
			}
		}
		return ""
	}
}

func (network *Network) SendStoreMessage(data []byte, contact *Contact) bool {
	netDialConnection, netDialError := net.Dial("udp4", contact.Address)
	if netDialError != nil {
		log.Println(netDialError)
		return false
	}
	defer func(netDialConnection net.Conn) {
		closeError := netDialConnection.Close()
		if closeError != nil {

		}
	}(netDialConnection)
	marshalBody, _ := json.Marshal(data)
	message := []byte("StoreMessage" + " " + string(marshalBody) + " ")
	newMarshalBody, _ := json.Marshal(network.CurrentNode)
	message = append(message, newMarshalBody...)
	_, writeError := netDialConnection.Write(message)
	if writeError != nil {
		return false
	}
	buffer := make([]byte, 4096)
	setReadDeadlineError := netDialConnection.SetReadDeadline(time.Now().Add(2 * time.Second))
	if setReadDeadlineError != nil {
		return false
	}
	end, readError := netDialConnection.Read(buffer)
	if readError != nil {
		return false
	}
	handleStoreResponse(buffer[:end], network)
	return true
}

func handleStoreResponse(message []byte, network *Network) {
	if string(message[:5]) == "Error" {
		log.Println(string(message))
		return
	} else {
		var storeContact *Contact
		unmarshalError := json.Unmarshal(message, &storeContact)
		if unmarshalError != nil {
			return
		}
		if contactUsability(storeContact, network) {
			network.RoutingTable.AddContact(*storeContact)
		}
	}
}

func contactUsability(contact *Contact, network *Network) bool {
	if contact == nil || contact.Address == "" || contact.ID == nil || contact.ID.Equals(network.CurrentNode.ID) {
		return false
	}
	return true
}
