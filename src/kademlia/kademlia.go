package kademlia

type Kademlia struct {
	nt *Network
	rt *RoutingTable
	m  map[string][]byte
}

func newKadmlia(nt *Network, rt *RoutingTable) Kademlia {
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
