package kademlia

import (
	"testing"
	"time"


	"github.com/stretchr/testify/assert"
)

func TestLookup(t *testing.T) {

	contact1 := NewContact(NewKademliaID("0000000000000000000000000000000000000001"), "localhost:8000")
	contact2 := NewContact(NewKademliaID("0000000000000000000000000000000000000002"), "localhost:8001")
	contact3 := NewContact(NewKademliaID("0000000000000000000000000000000000000003"), "localhost:8002")
	contact4 := NewContact(NewKademliaID("E000000000000000000000000000000000000001"), "localhost:8003")
	contact5 := NewContact(NewKademliaID("EE00000000000000000000000000000000000001"), "localhost:8004")

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

	kademlia1.Network.RoutingTable.AddContact(contact2)
	kademlia2.Network.RoutingTable.AddContact(contact3)
	kademlia3.Network.RoutingTable.AddContact(contact4)
	kademlia4.Network.RoutingTable.AddContact(contact5)
	kademlia3.Network.RoutingTable.AddContact(contact5)
	kademlia4.Network.RoutingTable.AddContact(contact2)
	kademlia3.Network.RoutingTable.AddContact(contact1)

	contact1.CalcDistance(contact1.ID)
	contact2.CalcDistance(contact1.ID)
	contact3.CalcDistance(contact1.ID)
	contact4.CalcDistance(contact1.ID)
	contact5.CalcDistance(contact1.ID)

	wantedList := []KademliaID{*contact1.ID, *contact3.ID, *contact2.ID}

	go Listen("localhost", "8000", &kademlia1)
	go Listen("localhost", "8001", &kademlia2)
	go Listen("localhost", "8002", &kademlia3)
	go Listen("localhost", "8003", &kademlia4)
	go Listen("localhost", "8004", &kademlia5)

	time.Sleep(1 * time.Second)

	contacts := kademlia1.LookupContact(contact1)
	var idList []KademliaID
	for _, contact := range contacts {
		idList = append(idList, *contact.ID)
	}
	// contacts := kademlia1.Network.SendFindContactMessage(&contact3, contact5.ID)

	assert.Equal(t, wantedList, idList)
}
