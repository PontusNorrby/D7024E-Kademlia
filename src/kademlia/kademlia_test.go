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
