// https://www.geeksforgeeks.org/overview-of-testing-package-in-golang/
package kademlia

import (
	"fmt"
	"testing"
	"time"
)

var tempInt int = 0

func TestKademliaID_Less(t *testing.T) {
	testNode := ToKademliaID("0000000000000000000000000000000000000000")
	testNode1 := ToKademliaID("0000000000000000000000000000000000000001")
	lessRes := testNode.Less(testNode1)
	if !lessRes {
		t.Fail()
	}
}

func TestKademliaID_Equals(t *testing.T) {
	testNode := ToKademliaID("A000000000000000000000000000000000000000")
	testNode1 := ToKademliaID("A000000000000000000000000000000000000000")
	testResult := testNode.Equals(testNode1)
	if !testResult {
		t.Fail()
	}
}

func TestKademliaID_CalcDistance(t *testing.T) {
	testNode := ToKademliaID("0000000000000000000000000000000000000010")
	testNode1 := ToKademliaID("0000000000000000000000000000000000000001")
	distanceRes := testNode.CalcDistance(testNode1)
	//fmt.Println(distanceRes)
	if distanceRes.String() != "0000000000000000000000000000000000000011" {
		t.Fail()
	}
}

func TestToKademliaID(t *testing.T) {
	testNode := ToKademliaID("A0000000000000000000000000000000000000000")
	if testNode != nil {
		t.Fail()
	}
}

func TestKademliaId(t *testing.T) {
	testNode := ToKademliaID("A000000000000000000000000000000000000000")
	if testNode.String() != "a000000000000000000000000000000000000000" {
		t.Fail()
	}
}

func TestKademlia_LookupData(t *testing.T) {
	testContact := NewContact(NewRandomKademliaID(), "127.0.0.1")
	testKademlia := NewKademliaStruct(NewNetwork(&testContact))
	testToken := []byte("AA")
	fmt.Println(testToken)
	testHash := testKademlia.store(testToken)
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

func TestStore(t *testing.T) {
	testNodes := setupKademliaNodes(tempInt)
	tempInt = tempInt + 4
	testResult, _ := testNodes[3].StoreValue([]byte("SavedValue"))
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
}

func setupKademliaNodes(i int) []*Kademlia {
	testNode1 := ToKademliaID("AAAA000000000000000000000000000000000000")
	testNode2 := ToKademliaID("BBBB000000000000000000000000000000000000")
	testNode3 := ToKademliaID("CCCC000000000000000000000000000000000000")
	testNode4 := ToKademliaID("DDDD000000000000000000000000000000000000")
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
	fmt.Println("Contact 1", testKademlia.Network.CurrentNode.ID)
	fmt.Println("Contact 2", testKademlia2.Network.CurrentNode.ID)
	fmt.Println("Contact 3", testKademlia3.Network.CurrentNode.ID)
	fmt.Println("Contact 4", testKademlia4.Network.CurrentNode.ID)
	time.Sleep(1 * time.Second)
	testArray := make([]*Kademlia, 4)
	testArray[0] = testKademlia
	testArray[1] = testKademlia2
	testArray[2] = testKademlia3
	testArray[3] = testKademlia4
	return testArray
}
