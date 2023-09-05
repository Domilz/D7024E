package kademlia

type Kademlia struct {
	RoutingTable *RoutingTable
	Self         *Contact
}

func NewKademlia() *Kademlia {
	contact := NewContact(NewRandomKademliaID(), "localhost:8080")
	return &Kademlia{
		RoutingTable: NewRoutingTable(contact),
		Self:         &contact,
	}
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
