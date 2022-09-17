package kademlia

type Network struct {
	Kademlia     *Kademlia
	MyNode       *Contact
	RoutingTable *RoutingTable
}

func NewNetwork(Kademlia *Kademlia, MyNode *Contact, rt *RoutingTable) Network {
	network := Network{}
	network.Kademlia = Kademlia
	network.RoutingTable = rt
	network.MyNode = MyNode
	return network
}

func Listen(ip string, port int) {
	// TODO
}

func (network *Network) SendPingMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}
