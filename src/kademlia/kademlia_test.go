// https://www.geeksforgeeks.org/overview-of-testing-package-in-golang/
package kademlia

import (
	"fmt"
	"testing"
)

func TestKademliaID_Less(t *testing.T) {
	idNode1 := NewKademliaID("0000000000000000000000000000000000000000")
	idNode2 := NewKademliaID("0000000000000000000000000000000000000001")
	lessRes := idNode1.Less(idNode2)
	if !lessRes {
		t.Fail()
	}
}

func TestKademliaID_Equals(t *testing.T) {
	idNode1 := NewKademliaID("0000000000000000000000000000000000000000")
	idNode2 := NewKademliaID("0000000000000000000000000000000000000000")
	equalsRes := idNode1.Equals(idNode2)
	if !equalsRes {
		t.Fail()
	}
}

func TestKademliaID_CalcDistance(t *testing.T) {
	idNode1 := NewKademliaID("0000000000000000000000000000000000000010")
	idNode2 := NewKademliaID("0000000000000000000000000000000000000001")
	distanceRes := idNode1.CalcDistance(idNode2)
	fmt.Println(distanceRes)
	if distanceRes.String() != "0000000000000000000000000000000000000011" {
		t.Fail()
	}
}

func TestGetValueSelf(t *testing.T) {
	testContact := NewContact(NewRandomKademliaID(), "172.20.0.1")
	testKademlia := NewKademliaStruct(NewNetwork(&testContact))
	testValue := []byte("FF00")
	varStore := testKademlia.Store(testValue)
	varTestValue, testContact := testKademlia.getData(&varStore)
	if varTestValue == nil {
		t.Fail()
	}
}

// Sets up 3 kademlia nodes
func setupKademlia() []*Kademlia {
	node1 := NewKademliaID("FFFFFFFF00000000000000000000000000000000")
	node2 := NewKademliaID("FFFFFFF000000000000000000000000000000000")
	node3 := NewKademliaID("FFFFFF0000000000000000000000000000000000")
	contact1 := NewContact(node1, "172.20.0.1")
	contact2 := NewContact(node2, "172.20.0.2")
	contact3 := NewContact(node3, "172.20.0.3")
	network1 := NewNetwork(&contact1)
	network2 := NewNetwork(&contact2)
	network3 := NewNetwork(&contact3)
	kademlia1 := NewKademliaStruct(network1)
	kademlia2 := NewKademliaStruct(network2)
	kademlia3 := NewKademliaStruct(network3)

	go kademlia1.Network.Listen("172.20.0.1", 1000, kademlia1)
	go kademlia2.Network.Listen("172.20.0.2", 1001, kademlia2)
	go kademlia3.Network.Listen("172.20.0.3", 1002, kademlia3)
	//kademlia2 knows of 1 and is put into kademlia 3 who is put into kademlia4
	kademlia2.Network.RoutingTable.AddContact(contact1)
	kademlia3.Network.RoutingTable.AddContact(contact2)
	arrKademlia := make([]*Kademlia, 3)
	arrKademlia[0] = kademlia1
	arrKademlia[1] = kademlia2
	arrKademlia[2] = kademlia3

	fmt.Println(arrKademlia)
	return arrKademlia
}
