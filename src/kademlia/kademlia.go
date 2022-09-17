package kademlia

import (
	"github.com/PontusNorrby/D7024E-Kademlia/src/network"
	"github.com/PontusNorrby/D7024E-Kademlia/src/routing"
)

type Kademlia struct {
	nt *network.Network
	rt *routing.RoutingTable
	m  map[string][]byte
}

func newKadmlia(nt *network.Network, rt *routing.RoutingTable) Kademlia {
	kademlia := Kademlia{}
	kademlia.nt = nt
	kademlia.rt = rt
	kademlia.m = make(map[string][]byte)

	return kademlia
}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	// TODO
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
