package kademlia

type Kademlia struct {
	m       map[KademliaID]Value
	Network *Network
}

type Value struct {
	data []byte
}

func NewKademliaStruct(network *Network) *Kademlia {
	kademlia := &Kademlia{}
	kademlia.m = make(map[KademliaID]Value)
	kademlia.Network = network
	return kademlia
}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	// TODO
	//iterativ find_node

}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) KademliaID {
	storeId := NewKademliaID(string(data))
	dataStore := Value{data: data}
	kademlia.m[*storeId] = dataStore
	return KademliaID{}
}
