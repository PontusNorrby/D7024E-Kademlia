package network

import (
	"github.com/PontusNorrby/D7024E-Kademlia/src/kademlia"
	"github.com/PontusNorrby/D7024E-Kademlia/src/routing"
)

type Network struct {
	Kademlia     *kademlia.Kademlia
	MyNode       *kademlia.Contact
	RoutingTable *routing.RoutingTable
}

func NewNetwork(Kademlia *kademlia.Kademlia, MyNode *kademlia.Contact, rt *routing.RoutingTable) Network {
	network := Network{}
	network.Kademlia = Kademlia
	network.RoutingTable = rt
	network.MyNode = MyNode
	return network
}

func Listen(ip string, port int) {
	// TODO
}

func (network *Network) SendPingMessage(contact *kademlia.Contact) {
	// TODO
}

func (network *Network) SendFindContactMessage(contact *kademlia.Contact) {
	// TODO
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}
