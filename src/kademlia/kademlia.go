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

func (kademlia *Kademlia) LookupContact(target *KademliaID) ContactCandidates {
	closerContacts := kademlia.Network.RoutingTable.FindClosestContacts(target, bucketSize)
	contactCandidates := kademlia.lookupContactTest(target, closerContacts)
	nodeContacts := contactCandidates.contacts
	if target.Equals(kademlia.Network.RoutingTable.me.ID) || 20 > len(nodeContacts) {
		contact := kademlia.Network.RoutingTable.me
		contact.CalcDistance(target)
		nodeContacts = append([]Contact{contact}, nodeContacts...)
		contactCandidates := ContactCandidates{nodeContacts}
		contactCandidates.Sort()
		return contactCandidates
	}
	return contactCandidates
}

func (kademlia *Kademlia) lookupContactTest(target *KademliaID, earlierContacts []Contact) ContactCandidates {
	foundContacts := 0
	newRouteFinding := NewRoutingTable(*kademlia.Network.CurrentNode)
	minLength := minVal(alphaValue, len(earlierContacts))
	for i := 0; i < minLength; i++ {
		newContact := earlierContacts[i]
		go func(newContact Contact) {
			contactsFetched := kademlia.Network.SendFindContactMessage(&newContact, target)
			for _, middleContact := range contactsFetched {
				newRouteFinding.AddContact(middleContact)
			}
		}(newContact)
	}
	contactClosest := newRouteFinding.FindClosestContacts(target, bucketSize)
	for _, contact := range contactClosest {
		for _, previousContact := range earlierContacts {
			if contact.ID.Equals(previousContact.ID) {
				foundContacts += 1
				break
			}
		}
	}
	if foundContacts == len(contactClosest) {
		return ContactCandidates{contactClosest}
	} else {
		return kademlia.lookupContactTest(target, contactClosest)
	}
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) storeHelp(data []byte) ([]*KademliaID, string) {
	target := NewKademliaID(string(data))
	closest := kademlia.LookupContact(target)
	var storedNodes []*KademliaID
	for _, contact := range closest.contacts {
		if contact.ID.Equals(kademlia.Network.RoutingTable.me.ID) {
			kademlia.Store(data)
			storedNodes = append(storedNodes, contact.ID)
		}
		go func(contact Contact) {
			res := kademlia.Network.SendStoreMessage(data, &contact, kademlia)
			if res {

			}
		}(contact)
	}
	return storedNodes, target.String()
}

// Just stores the data in this node not on the "correct" node
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
