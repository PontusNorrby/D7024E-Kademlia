package cli

import (
	"bufio"
	"fmt"
	"github.com/PontusNorrby/D7024E-Kademlia/src/kademlia"
	. "github.com/PontusNorrby/D7024E-Kademlia/src/kademlia"
	"os"
	"strings"
)

func StartCLI(nodeShutDown func(), kademlia *kademlia.Kademlia) {
	run(userInput, nodeShutDown, kademlia)
}

func run(userInput func() string, nodeShutDown func(), kademlia *kademlia.Kademlia) {
	for {
		input := userInput()
		if input == "exit" {
			fmt.Println("You are going to exit, to confirm please write yes else, press enter")
			confirmation := userInput()
			if confirmation == "yes" {
				nodeShutDown()
				return
			}
		} else if input == "put" {
			store(userInput, kademlia.StoreHelp)
		} else if input == "get" {
			get(userInput, kademlia)
		} else if input == "help" {
			helpList()
		} else {
			fmt.Println("Unknown command, please write help to get a list of commands")
		}

	}
}

func store(input func() string, Store func(data []byte) ([]*KademliaID, string)) {
	fmt.Println("What value do you want to store? ")
	value := input()
	storeID, hash := Store([]byte(value))
	fmt.Println("Hash of", value, "is", hash)
	fmt.Println("the value stored in nodes: ", storeID)
}

func get(input func() string, kademlia *kademlia.Kademlia) {
	fmt.Println("insert the ID you want to get")
	stringValue := input()
	id := ToKademliaID(stringValue)
	value, contact := kademlia.GetData(id)
	if value == nil {
		fmt.Println("value not found")
		return
	}
	fmt.Println("\""+*value+"\" found at node", contact.ID)
}

//func store(input func() string, store func(data []byte) kademlia.KademliaID) {

//}

func helpList() {
	fmt.Println("put - Takes a single argument, the contents of the file you are uploading, and outputs the hash of the object, if it could be uploaded successfully.")
	fmt.Println("get - Takes a hash as its only argument, and outputs the contents of the object and the node it was retrieved from, if it could be downloaded successfully.")
	fmt.Println("exit - Terminates the node.")
	fmt.Println("help - Prints this list")
}

func userInput() string {
	input := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	input.Scan()
	newInput := strings.Trim(input.Text(), " ")
	newInput = strings.ToLower(newInput)
	return newInput
}
