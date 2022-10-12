package kademlia

import (
	"testing"
	"time"
)

// To run tests:
// go test -coverprofile cover.out =./... ./...
// To see coverage:
// go tool cover -html cover.out

func TestPingNode(t *testing.T) {
	contact := NewContact(NewRandomKademliaID(), "127.0.0.1:8000")
	kademlia := NewKademliaStruct(NewNetwork(&contact))

	go kademlia.Network.Listen("127.0.0.1", 8000, kademlia)

	go kademlia.Network.SendPingMessage(&contact)

	time.Sleep(1 * time.Millisecond)

	return
}

func TestPingNode2(t *testing.T) {
	nodeID := NewRandomKademliaID()
	contact := NewContact(nodeID, "127.0.0.1:8000")
	network := NewNetwork(&contact)

	if network.SendPingMessage(&contact) {
		t.Fail()
	}
}
