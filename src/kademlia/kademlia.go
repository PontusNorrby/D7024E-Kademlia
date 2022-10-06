package kademlia

const alphaValue = 3

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

func (kademlia *Kademlia) LookupContact(target *KademliaID) {
	// TODO
	closerContacts := kademlia.Network.RoutingTable.FindClosestContacts(target, bucketSize)
	newRouteFinding := NewRoutingTable(*kademlia.Network.CurrentNode)
	minLength := minVal(alphaValue, len(closerContacts))
	for i := 0; i < minLength; i++ {
		newContact := closerContacts[i]
		contactsFetched := kademlia.Network.SendFindContactMessage(&newContact, target)
		for _, tempContact := range contactsFetched {
			newRouteFinding.AddContact(tempContact)
		}
	}

}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) KademliaID {
	storeId := NewKademliaID(string(data))
	dataStore := Value{data: data}
	kademlia.m[*storeId] = dataStore
	return *storeId
}

func minVal(x int, y int) int {
	if x <= y {
		return x
	}
	return y
}
