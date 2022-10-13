// https://www.geeksforgeeks.org/overview-of-testing-package-in-golang/
package kademlia

import (
	"fmt"
	"testing"
	"time"
)

var i int = 0

/*func TestKademliaID_Less(t *testing.T) {
	idNode1 := ToKademliaID("0000000000000000000000000000000000000000")
	idNode2 := ToKademliaID("0000000000000000000000000000000000000001")
	lessRes := idNode1.Less(idNode2)
	if !lessRes {
		t.Fail()
	}
}

func TestKademliaID_Equals(t *testing.T) {
	idNode1 := ToKademliaID("A000000000000000000000000000000000000000")
	idNode2 := ToKademliaID("A000000000000000000000000000000000000000")
	equalsRes := idNode1.Equals(idNode2)
	if !equalsRes {
		t.Fail()
	}
}

func TestKademliaID_CalcDistance(t *testing.T) {
	idNode1 := ToKademliaID("0000000000000000000000000000000000000010")
	idNode2 := ToKademliaID("0000000000000000000000000000000000000001")
	distanceRes := idNode1.CalcDistance(idNode2)
	//fmt.Println(distanceRes)
	if distanceRes.String() != "0000000000000000000000000000000000000011" {
		t.Fail()
	}
}*/

func TestFindContact(t *testing.T) {
	kademliaNodes := returnKademliaNodes(i)
	i = i + 4
	res := kademliaNodes[3].LookupContact(kademliaNodes[0].Network.CurrentNode.ID).contacts
	fmt.Println(res, kademliaNodes[0].Network.CurrentNode.ID)
	if res == nil || !res[0].ID.Equals(kademliaNodes[0].Network.CurrentNode.ID) {
		t.Fail()
	}
}

/*
	func TestFindContact2(t *testing.T) {
		kademliaNodes := returnKademliaNodes(i)
		i = i + 4
		res := kademliaNodes[3].LookupContact(kademliaNodes[1].Network.CurrentNode.ID).contacts
		if res == nil || !res[0].ID.Equals(kademliaNodes[1].Network.CurrentNode.ID) {
			t.Fail()
		}
	}

	func TestFindContact3(t *testing.T) {
		kademliaNodes := returnKademliaNodes(i)
		i = i + 4
		res := kademliaNodes[3].LookupContact(kademliaNodes[2].Network.CurrentNode.ID).contacts
		if res == nil || !res[0].ID.Equals(kademliaNodes[2].Network.CurrentNode.ID) {
			t.Fail()
		}
	}

	func TestFindContact4(t *testing.T) {
		kademliaNodes := returnKademliaNodes(i)
		i = i + 4
		kadId := ToKademliaID(kademliaNodes[3].Network.CurrentNode.ID.String())
		res := kademliaNodes[3].LookupContact(kadId).contacts
		fmt.Println(res, kademliaNodes[3].Network.CurrentNode)
		if res == nil || !res[0].ID.Equals(kademliaNodes[3].Network.CurrentNode.ID) {
			t.Fail()
		}
	}
*/
func returnKademliaNodes(i int) []*Kademlia {
	nodeID := ToKademliaID("A000000000000000000000000000000000000000")
	contact := NewContact(nodeID, ("127.0.0.1:" + fmt.Sprint(7000+i)))
	network := NewNetwork(&contact)
	kademlia := NewKademliaStruct(network)
	nodeID2 := ToKademliaID("B000000000000000000000000000000000000000")
	contact2 := NewContact(nodeID2, ("127.0.0.1:" + fmt.Sprint(7001+i)))
	network2 := NewNetwork(&contact2)
	kademlia2 := NewKademliaStruct(network2)
	nodeID3 := ToKademliaID("C000000000000000000000000000000000000000")
	contact3 := NewContact(nodeID3, ("127.0.0.1:" + fmt.Sprint(7002+i)))
	network3 := NewNetwork(&contact3)
	kademlia3 := NewKademliaStruct(network3)
	nodeID4 := ToKademliaID("D000000000000000000000000000000000000000")
	contact4 := NewContact(nodeID4, ("127.0.0.1:" + fmt.Sprint(7003+i)))
	network4 := NewNetwork(&contact4)
	kademlia4 := NewKademliaStruct(network4)

	go kademlia.Network.Listen("127.0.0.1", 7000+i, kademlia)
	go kademlia2.Network.Listen("127.0.0.1", 7001+i, kademlia2)
	go kademlia3.Network.Listen("127.0.0.1", 7002+i, kademlia3)
	go kademlia3.Network.Listen("127.0.0.1", 7003+i, kademlia4)
	kademlia2.Network.RoutingTable.AddContact(contact)
	kademlia3.Network.RoutingTable.AddContact(contact2)
	kademlia4.Network.RoutingTable.AddContact(contact3)
	fmt.Println("Contact 1", kademlia.Network.CurrentNode.ID)
	fmt.Println("Contact 2", kademlia2.Network.CurrentNode.ID)
	fmt.Println("Contact 3", kademlia3.Network.CurrentNode.ID)
	fmt.Println("Contact 4", kademlia4.Network.CurrentNode.ID)
	time.Sleep(1 * time.Second)
	kademliaArray := make([]*Kademlia, 4)
	kademliaArray[0] = kademlia
	kademliaArray[1] = kademlia2
	kademliaArray[2] = kademlia3
	kademliaArray[3] = kademlia4
	return kademliaArray
}

/*Sets up 3 kademlia nodes
func setupKademlia() []*Kademlia {
	node1 := ToKademliaID("FFFFFFFF00000000000000000000000000000000")
	node2 := ToKademliaID("FFFFFFF000000000000000000000000000000000")
	node3 := ToKademliaID("FFFFFF0000000000000000000000000000000000")
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
	fmt.Println("Contact 1: ", kademlia1.Network.CurrentNode.ID)
	fmt.Println("Contact 2: ", kademlia2.Network.CurrentNode.ID)
	fmt.Println("Contact 3: ", kademlia3.Network.CurrentNode.ID)
	time.Sleep(1 * time.Second)
	arrKademlia := make([]*Kademlia, 3)
	arrKademlia[0] = kademlia1
	arrKademlia[1] = kademlia2
	arrKademlia[2] = kademlia3

	return arrKademlia
}*/
