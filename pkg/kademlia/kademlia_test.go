package kademlia

import (
	"fmt"
	"testing"
)

func TestLookup(t *testing.T) {

	contact1 := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000")
	contact2 := NewContact(NewKademliaID("EFFFFFFF00000000000000000000000000000000"), "localhost:8001")
	contact3 := NewContact(NewKademliaID("1111111200000000000000000000000000000000"), "localhost:8002")
	contact4 := NewContact(NewKademliaID("1111111400000000000000000000000000000000"), "localhost:8003")
	contact5 := NewContact(NewKademliaID("1111111500000000000000000000000000000000"), "localhost:8004")

	kademlia1 := Kademlia{
		Network: &Network{
			RoutingTable: NewRoutingTable(contact1),
			Self:         &contact1,
		},
	}
	kademlia2 := Kademlia{
		Network: &Network{
			RoutingTable: NewRoutingTable(contact2),
			Self:         &contact2,
		},
	}
	kademlia3 := Kademlia{
		Network: &Network{
			RoutingTable: NewRoutingTable(contact3),
			Self:         &contact3,
		},
	}
	kademlia4 := Kademlia{
		Network: &Network{
			RoutingTable: NewRoutingTable(contact4),
			Self:         &contact4,
		},
	}
	kademlia5 := Kademlia{
		Network: &Network{
			RoutingTable: NewRoutingTable(contact5),
			Self:         &contact5,
		},
	}
	fmt.Println("första")

	kademlia1.Network.RoutingTable.AddContact(contact2)
	kademlia1.Network.RoutingTable.AddContact(contact3)
	kademlia1.Network.RoutingTable.AddContact(contact4)
	kademlia2.Network.RoutingTable.AddContact(contact4)
	kademlia3.Network.RoutingTable.AddContact(contact5)
	kademlia4.Network.RoutingTable.AddContact(contact5)

	go Listen("localhost", "8000", &kademlia1)
	go Listen("localhost", "8001", &kademlia2)
	go Listen("localhost", "8002", &kademlia3)
	go Listen("localhost", "8003", &kademlia4)
	go Listen("localhost", "8004", &kademlia5)
	fmt.Println("andra")

	contacts := kademlia1.LookupContact(contact5)
	// contacts := kademlia1.Network.SendFindContactMessage(&contact3, contact5.ID)

	for _, c := range contacts {
		fmt.Println("Det vi får ut", c)
	}
}
