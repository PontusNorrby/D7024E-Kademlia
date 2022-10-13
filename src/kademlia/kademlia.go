package kademlia

import (
	"fmt"
	"sync"
)

const alphaValue = 3

type Kademlia struct {
	m            map[KademliaID]*Value
	Network      *Network
	KnownHolders map[Contact]KademliaID
}

type Value struct {
	data []byte
}

func NewKademliaStruct(network *Network) *Kademlia {
	kademlia := &Kademlia{}
	kademlia.m = make(map[KademliaID]*Value)
	kademlia.Network = network
	kademlia.KnownHolders = make(map[Contact]KademliaID)
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

func (kademlia *Kademlia) LookupData(valueSelf KademliaID) []byte {
	value, check := kademlia.m[valueSelf]
	if check {
		return value.data
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
		var resultString *string = nil
		var contactCandidates Contact = Contact{}
		for i := 0; i < length; i++ {
			go func(possibleContact Contact) {
				findDataRes := kademlia.Network.SendFindDataMessage(value, &possibleContact)
				if !(findDataRes == "Error") {
					resultString = &findDataRes
					contactCandidates = possibleContact
				}
			}(possibleContacts[0])
			possibleContacts = possibleContacts[1:]
		}
		if resultString != nil {
			return resultString, contactCandidates
		}
	}
	return nil, Contact{}
}

func (kademlia *Kademlia) StoreValue(data []byte) ([]*KademliaID, string) {
	target := NewKademliaID(string(data))
	closest := kademlia.LookupContact(target)
	var storedNodes []*KademliaID
	for _, contact := range closest.contacts {
		if contact.ID.Equals(kademlia.Network.RoutingTable.me.ID) {
			kademlia.store(data)
			storedNodes = append(storedNodes, contact.ID)
		}
		go func(contact Contact) {
			res := kademlia.Network.SendStoreMessage(data, &contact, kademlia)
			if res {
				storedNodes = append(storedNodes, contact.ID)
			}
		}(contact)
	}
	return storedNodes, target.String()
}

// Just stores the data in this node not on the "correct" node
func (kademlia *Kademlia) store(data []byte) KademliaID {
	fmt.Println("3")
	storeId := NewKademliaID(string(data))
	fmt.Println("4")
	dataStore := Value{data}
	fmt.Println("5")
	kademlia.m[*storeId] = &dataStore
	fmt.Println("6")
	return *storeId
}

func (kademlia *Kademlia) AddToKnown(contact *Contact, hash *KademliaID) {
	kademlia.KnownHolders[*contact] = *hash
}

func minVal(x int, y int) int {
	if x <= y {
		return x
	}
	return y
}
