package kademlia

type Ping struct {
	startMessage string
}

func newPing() *Ping {
	return &Ping{"Ping"}
}

type FindContact struct {
	startMessage string
}

func newFindContact() *FindContact {
	return &FindContact{"FindContact"}
}

type FindData struct {
	startMessage string
}

func newFindData() *FindData {
	return &FindData{"FindData"}
}

type StoreMessage struct {
	startMessage string
}

func newStoreMessage() *StoreMessage {
	return &StoreMessage{"StoreMessage"}
}
