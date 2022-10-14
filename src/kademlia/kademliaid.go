package kademlia

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"math/rand"
)

// the static number of bytes in a KademliaID
const IDLength = 20

// type definition of a KademliaID
type KademliaID [IDLength]byte

// NewKademliaID returns a new instance of a KademliaID based on the string input
/*func NewKademliaID(data string) *KademliaID {
	decoded, _ := hex.DecodeString(data)
	newDecoded := hex.EncodeToString(decoded[0:IDLength])

	newKademliaID := KademliaID{}
	for i := 0; i < IDLength; i++ {
		newKademliaID[i] = newDecoded[i]
	}

	return &newKademliaID
}*/

func NewKademliaID(data string) *KademliaID {
	hashBytes := sha1.Sum([]byte(data))
	hash := hex.EncodeToString(hashBytes[0:IDLength])

	newKademliaID := KademliaID{}
	for i := 0; i < IDLength; i++ {
		newKademliaID[i] = hash[i]
	}

	return &newKademliaID
}

func ToKademliaID(data string) *KademliaID {
	if len(data) < 40 {
		return nil
	}
	res, err := hex.DecodeString(data)
	if err != nil {
		fmt.Println("FAILED TO DECODE KADEMLIA ID", err)
		return nil
	} else {
		return (*KademliaID)(res)
	}
}

// NewRandomKademliaID returns a new instance of a random KademliaID,
// change this to a better version if you like
func NewRandomKademliaID() *KademliaID {
	newKademliaID := KademliaID{}
	for i := 0; i < IDLength; i++ {
		newKademliaID[i] = uint8(rand.Intn(256))
	}
	return &newKademliaID
}

// Less returns true if kademliaID < otherKademliaID (bitwise)
func (kademliaID KademliaID) Less(otherKademliaID *KademliaID) bool {
	for i := 0; i < IDLength; i++ {
		if kademliaID[i] != otherKademliaID[i] {
			return kademliaID[i] < otherKademliaID[i]
		}
	}
	return false
}

// Equals returns true if kademliaID == otherKademliaID (bitwise)
func (kademliaID KademliaID) Equals(otherKademliaID *KademliaID) bool {
	for i := 0; i < IDLength; i++ {
		if kademliaID[i] != otherKademliaID[i] {
			return false
		}
	}
	return true
}

// CalcDistance returns a new instance of a KademliaID that is built
// through a bitwise XOR operation betweeen kademliaID and target
func (kademliaID KademliaID) CalcDistance(target *KademliaID) *KademliaID {
	result := KademliaID{}
	for i := 0; i < IDLength; i++ {
		result[i] = kademliaID[i] ^ target[i]
	}
	return &result
}

// String returns a simple string representation of a KademliaID
func (kademliaID *KademliaID) String() string {
	return hex.EncodeToString(kademliaID[0:IDLength])
}
