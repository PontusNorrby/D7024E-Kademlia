package cli

import (
	"bufio"
	"fmt"
	"github.com/PontusNorrby/D7024E-Kademlia/src/kademlia"
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
			fmt.Println("You are going to exit, to confirm please write yes >")
			confirmation := userInput()
			if confirmation == "yes" {
				nodeShutDown()
				return
			}
		} else if input == "put" {
			store(userInput, kademlia.Store)
		} else if input == "get" {
			get(userInput, kademlia)
		} else if input == "help" {
			helpList()
		} else {
			fmt.Println("Unknown command, please write help to get a list of commands")
		}

	}
}

// TODO
func get(userInput func() string, k *kademlia.Kademlia) {
	fmt.Println("please insert the value you want to get")

}

func store(userInput func() string, store func(data []byte) kademlia.KademliaID) {
	fmt.Println("please insert the value you want to store")
	value := userInput()
	fmt.Println("Stored in nodes: ", store([]byte(value)))

}

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
