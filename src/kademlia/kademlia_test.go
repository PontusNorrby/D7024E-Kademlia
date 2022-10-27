// https://www.geeksforgeeks.org/overview-of-testing-package-in-golang/
package kademlia

import (
	"fmt"
	"testing"
	"time"
)

var tempInt = 0

func TestKademliaID_Less(t *testing.T) {
	testNode := IDDecoder("0000000000000000000000000000000000000000")
	testNode1 := IDDecoder("0000000000000000000000000000000000000001")
	lessRes := testNode.Less(testNode1)
	if !lessRes {
		t.Fail()
	}
}

func TestKademliaID_Equals(t *testing.T) {
	testNode := IDDecoder("A000000000000000000000000000000000000000")
	testNode1 := IDDecoder("A000000000000000000000000000000000000000")
	testResult := testNode.Equals(testNode1)
	if !testResult {
		t.Fail()
	}
}

func TestKademliaID_CalcDistance(t *testing.T) {
	testNode := IDDecoder("0000000000000000000000000000000000000010")
	testNode1 := IDDecoder("0000000000000000000000000000000000000001")
	distanceRes := testNode.CalcDistance(testNode1)
	if distanceRes.String() != "0000000000000000000000000000000000000011" {
		t.Fail()
	}
}

func TestToKademliaID(t *testing.T) {
	testNode := IDDecoder("A0000000000000000000000000000000000000000")
	if testNode != nil {
		t.Fail()
	}
}

func TestKademliaId(t *testing.T) {
	testNode := IDDecoder("A000000000000000000000000000000000000000")
	if testNode.String() != "a000000000000000000000000000000000000000" {
		t.Fail()
	}
}

func TestKademlia_LookupData(t *testing.T) {
	testContact := NewContact(NewRandomKademliaID(), "127.0.0.1")
	testKademlia := NewKademliaStruct(NewNetwork(&testContact))
	testToken := []byte("AA")
	fmt.Println(testToken)
	testHash := testKademlia.storeDataHelp(testToken)
	response := testKademlia.LookupData(testHash)
	if response == nil {
		t.Fail()
	}
	response = testKademlia.LookupData(*NewRandomKademliaID())

	if response != nil {
		t.Fail()
	}
}

func TestFindContact(t *testing.T) {
	testNodes := setupKademliaNodes(tempInt)
	tempInt = tempInt + 4
	testResult := testNodes[3].LookupContact(testNodes[0].Network.CurrentNode.ID).contacts

	if testResult == nil || !testResult[0].ID.Equals(testNodes[0].Network.CurrentNode.ID) {
		t.Fail()
	}
}

func TestStoreData(t *testing.T) {
	testNodes := setupKademliaNodes(tempInt)
	tempInt = tempInt + 4
	testResult, _ := testNodes[3].StoreData([]byte("SavedValue"))
	if len(testResult) != len(testNodes) {
		fmt.Println("STORED ON:", testResult, "\nTOTAL NODES:", len(testNodes))
		t.FailNow()
	}
	for i, node := range testNodes {
		if len(node.m) == 0 {
			fmt.Println("Node", i, " have", len(node.m))
			t.FailNow()
		}
	}
	fmt.Println("Cleared Store test")
}

func TestStoreDataAndFindData(t *testing.T) {
	nodeID := NewRandomKademliaID()
	nodeID2 := NewRandomKademliaID()
	contact := NewContact(nodeID, "127.0.0.1:8002")
	contact2 := NewContact(nodeID2, "127.0.0.1:8003")
	network := NewNetwork(&contact)
	network2 := NewNetwork(&contact2)
	kademlia := NewKademliaStruct(network)
	kademlia2 := NewKademliaStruct(network2)

	go network.Listen("127.0.0.1", 8002, kademlia)
	go network2.Listen("127.0.0.1", 8003, kademlia2)

	time.Sleep(1 * time.Millisecond)
	fmt.Println("Testline1")
	fmt.Println(network2.SendStoreMessage([]byte("String"), &contact))

	time.Sleep(1 * time.Millisecond)

	hash := NewKademliaID("String")
	res := network2.SendFindDataMessage(hash, &contact)
	if res != "String" {
		fmt.Println("Res is", res)
		fmt.Println("Broken here")
		t.Fail()
	}
	time.Sleep(6 * time.Second)

	return
}

func setupKademliaNodes(i int) []*Kademlia {
	testNode1 := IDDecoder("AAAA000000000000000000000000000000000000")
	testNode2 := IDDecoder("BBBB000000000000000000000000000000000000")
	testNode3 := IDDecoder("CCCC000000000000000000000000000000000000")
	testNode4 := IDDecoder("DDDD000000000000000000000000000000000000")
	testContact := NewContact(testNode1, ("127.0.0.1:" + fmt.Sprint(1000+i)))
	testContact2 := NewContact(testNode2, ("127.0.0.1:" + fmt.Sprint(1001+i)))
	testContact3 := NewContact(testNode3, ("127.0.0.1:" + fmt.Sprint(1002+i)))
	testContact4 := NewContact(testNode4, ("127.0.0.1:" + fmt.Sprint(1003+i)))
	testNetwork := NewNetwork(&testContact)
	testNetwork2 := NewNetwork(&testContact2)
	testNetwork3 := NewNetwork(&testContact3)
	testNetwork4 := NewNetwork(&testContact4)
	testKademlia := NewKademliaStruct(testNetwork)
	testKademlia2 := NewKademliaStruct(testNetwork2)
	testKademlia3 := NewKademliaStruct(testNetwork3)
	testKademlia4 := NewKademliaStruct(testNetwork4)

	go testKademlia.Network.Listen("127.0.0.1", 1000+i, testKademlia)
	go testKademlia2.Network.Listen("127.0.0.1", 1001+i, testKademlia2)
	go testKademlia3.Network.Listen("127.0.0.1", 1002+i, testKademlia3)
	go testKademlia3.Network.Listen("127.0.0.1", 1003+i, testKademlia4)
	testKademlia2.Network.RoutingTable.AddContact(testContact)
	testKademlia3.Network.RoutingTable.AddContact(testContact2)
	testKademlia4.Network.RoutingTable.AddContact(testContact3)
	time.Sleep(1 * time.Second)
	testArray := make([]*Kademlia, 4)
	testArray[0] = testKademlia
	testArray[1] = testKademlia2
	testArray[2] = testKademlia3
	testArray[3] = testKademlia4
	return testArray
}

func TestNetworkStruct(t *testing.T) {
	nodeID := NewRandomKademliaID()
	contact := NewContact(nodeID, "127.0.0.1:3000")

	network := NewNetwork(&contact)

	fmt.Println(network.CurrentNode.ID)
}

func TestPingNode(t *testing.T) {
	nodeID := NewRandomKademliaID()
	contact := NewContact(nodeID, "127.0.0.1:8000")
	network := NewNetwork(&contact)
	kademlia := NewKademliaStruct(network)

	go network.Listen("127.0.0.1", 8000, kademlia)

	go network.SendPingMessage(&contact)

	time.Sleep(1 * time.Millisecond)

	return
}

func TestPingFalseAddress(t *testing.T) {
	nodeID := NewRandomKademliaID()
	contact := NewContact(nodeID, "127.0.:8000")
	network := NewNetwork(&contact)

	if network.SendPingMessage(&contact) {
		t.Fail()
	}
}

func TestFindNode(t *testing.T) {
	nodeID := NewRandomKademliaID()
	contact := NewContact(nodeID, "127.0.0.1:8001")
	network := NewNetwork(&contact)
	kademlia := NewKademliaStruct(network)

	go network.Listen("127.0.0.1", 8001, kademlia)

	go network.SendFindContactMessage(&contact, nodeID)

	time.Sleep(1 * time.Millisecond)

	return
}
