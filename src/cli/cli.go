package cli

import (
	"bufio"
	"fmt"
	"github.com/PontusNorrby/D7024E-Kademlia/src/kademlia"
	. "github.com/PontusNorrby/D7024E-Kademlia/src/kademlia"
	"os"
	"strings"
)

func StartCLI(nodeExit func(), kademlia *kademlia.Kademlia) {
	run(userInput, nodeExit, kademlia)
}

func run(userInput func() string, nodeExit func(), kademlia *kademlia.Kademlia) {
	for {
		input := userInput()
		if input == "exit" {
			fmt.Println("Are you sure you want exit the node(y/n)?")
			confirmation := userInput()
			if confirmation == "y" {
				nodeExit()
				return
			} else {
				continue
			}
		} else if input == "put" {
			store(userInput, kademlia.StoreValue)
		} else if input == "get" {
			get(userInput, kademlia)
		} else if input == "help" {
			helpList()
		} else {
			fmt.Println("Unknown command, please write help to get a list of available commands")
		}

	}
}

func store(input func() string, Store func(data []byte) ([]*KademliaID, string)) {
	fmt.Println("Please insert the value you want to store...")
	value := input()
	storeID, target := Store([]byte(value))
	fmt.Println(
		"\bThe value you saved is                    	", value, "\n",
		"\bThe hash of the value is                      	", target, "\n",
		"\bYour value is saved in node(s) with id(:s)	", storeID)

}

func get(input func() string, kademlia *kademlia.Kademlia) {
	fmt.Println("Please insert the hash value to get the node(s) id(s)...")
	stringValue := input()
	id := ToKademliaID(stringValue)
	value, contact := kademlia.GetData(id)
	if value == nil {
		fmt.Println("No such hash value!")
		return
	}
	fmt.Println(*value+": found at node			", contact.ID)
}

func helpList() {
	fmt.Println("put - Store a value in the nodes, takes string and returns the hash value and the node id")
	fmt.Println("get - Get the stored value in the nodes, takes the hash value and returns the node id")
	fmt.Println("exit - Terminates the node.")
	fmt.Println("help - Prints help list")
}

func userInput() string {
	input := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	input.Scan()
	newInput := strings.Trim(input.Text(), " ")
	newInput = strings.ToLower(newInput)
	return newInput
}
