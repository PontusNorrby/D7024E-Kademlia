package kademlia

import (
	"sync"
)

const alphaValue = 3

type Kademlia struct {
	m          map[KademliaID]*Value
	storeMutex sync.Mutex
	Network    *Network
}

type Value struct {
	data []byte
}

func NewKademliaStruct(network *Network) *Kademlia {
	kademlia := &Kademlia{}
	kademlia.m = make(map[KademliaID]*Value)
	kademlia.Network = network
	kademlia.storeMutex = sync.Mutex{}
	return kademlia
}

func (kademlia *Kademlia) LookupContact(target *KademliaID) ContactCandidates {
	closerContacts := kademlia.Network.RoutingTable.FindClosestContacts(target, bucketSize)
	contactCandidatesS := kademlia.lookupContactHelp(target, closerContacts)
	nodeContacts := contactCandidatesS.contacts
	if target.Equals(kademlia.Network.RoutingTable.me.ID) || 20 > len(nodeContacts) {
		contact := kademlia.Network.RoutingTable.me
		contact.CalcDistance(target)
		nodeContacts = append([]Contact{contact}, nodeContacts...)
		contactCandidates := ContactCandidates{nodeContacts}
		contactCandidates.Sort()
		return contactCandidates
	}
	return contactCandidatesS
}

func (kademlia *Kademlia) lookupContactHelp(target *KademliaID, earlierContacts []Contact) ContactCandidates {
	newRouteFinding := NewRoutingTable(*kademlia.Network.CurrentNode)
	var goLock sync.WaitGroup
	minLength := minVal(alphaValue, len(earlierContacts))
	goLock.Add(minLength)
	for i := 0; i < minLength; i++ {
		newContact := earlierContacts[i]
		go func(newContact Contact) {
			defer goLock.Done()
			contactsFetched := kademlia.Network.SendFindContactMessage(&newContact, target)
			for _, middleContact := range contactsFetched {
				newRouteFinding.AddContact(middleContact)
			}
		}(newContact)
	}
	goLock.Wait()
	contactClosest := newRouteFinding.FindClosestContacts(target, bucketSize)
	foundContacts := 0
	for _, contact := range contactClosest {
		for _, previousContact := range earlierContacts {
			if contact.ID.Equals(previousContact.ID) {
				foundContacts++
				break
			}
		}
	}
	if foundContacts == len(contactClosest) {
		return ContactCandidates{contactClosest}
	} else {
		return kademlia.lookupContactHelp(target, contactClosest)
	}
}

func (kademlia *Kademlia) LookupData(value KademliaID) []byte {
	val, check := kademlia.m[value]
	if check {
		return val.data
	}
	return nil
}

func (kademlia *Kademlia) GetData(value *KademliaID) (*string, Contact) {
	selfCheck := kademlia.LookupData(*value)
	if selfCheck != nil {
		gottenValue := string(selfCheck)
		//If the function LookupData finds the data on the same node,
		//Return the value and the node that has the information (the node we are on)
		return &gottenValue, *kademlia.Network.CurrentNode
	}
	possibleContacts := kademlia.LookupContact(value).contacts
	for len(possibleContacts) > 0 {
		length := minVal(alphaValue, len(possibleContacts))
		var wg sync.WaitGroup
		wg.Add(length)
		var resultString *string = nil
		var contactCandidates = Contact{}
		for i := 0; i < length; i++ {
			go func(possibleContact Contact) {
				defer wg.Done()
				findDataRes := kademlia.Network.SendFindDataMessage(value, &possibleContact)
				if !(findDataRes == "Error") {
					resultString = &findDataRes
					contactCandidates = possibleContact
				}
			}(possibleContacts[0])
			possibleContacts = possibleContacts[1:]
		}
		wg.Wait()
		if resultString != nil {
			return resultString, contactCandidates
		}
	}
	return nil, Contact{}
}

func (kademlia *Kademlia) StoreData(data []byte) ([]*KademliaID, string) {
	target := NewKademliaID(string(data))
	closest := kademlia.LookupContact(target)
	var storedNodes []*KademliaID
	var wg sync.WaitGroup
	wg.Add(len(closest.contacts))
	for _, contact := range closest.contacts {
		if contact.ID.Equals(kademlia.Network.RoutingTable.me.ID) {
			kademlia.storeDataHelp(data)
			storedNodes = append(storedNodes, contact.ID)
			wg.Done()
			continue
		}
		go func(contact Contact) {
			defer wg.Done()
			res := kademlia.Network.SendStoreMessage(data, &contact)
			if res {
				storedNodes = append(storedNodes, contact.ID)
			}
		}(contact)
	}
	wg.Wait()
	return storedNodes, target.String()
}

// Just stores the data in this node not on the "correct" node
func (kademlia *Kademlia) storeDataHelp(data []byte) KademliaID {
	storeId := NewKademliaID(string(data))
	dataStore := Value{data}
	kademlia.m[*storeId] = &dataStore
	return *storeId
}

func minVal(x int, y int) int {
	if x <= y {
		return x
	}
	return y
}
