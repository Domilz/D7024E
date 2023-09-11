package kademlia

import (
	"fmt"
	"testing"
	"time"
)

func TestNetwork(t *testing.T) {

	contact1 := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000")
	contact2 := NewContact(NewKademliaID("EFFFFFFF00000000000000000000000000000000"), "localhost:8001")
	kademlia1 := Kademlia{
		Network: &Network{
			RoutingTable: NewRoutingTable(contact1),
			Self:         &contact1,
		},
	}

	kademlia1.Network.RoutingTable.AddContact(NewContact(NewKademliaID("1111111200000000000000000000000000000000"), "localhost:8002"))
	kademlia1.Network.RoutingTable.AddContact(NewContact(NewKademliaID("1111111300000000000000000000000000000000"), "localhost:8003"))
	kademlia1.Network.RoutingTable.AddContact(NewContact(NewKademliaID("1111111400000000000000000000000000000000"), "localhost:8004"))

	kademlia2 := Kademlia{
		Network: &Network{
			RoutingTable: NewRoutingTable(contact2),
			Self:         &contact2,
		},
	}

	kademlia2.Network.RoutingTable.AddContact(contact1)

	network2 := kademlia2.Network

	go Listen("localhost", "8000", &kademlia1)
	go Listen("localhost", "8001", &kademlia2)
	time.Sleep(5 * time.Second)
	target := NewContact(NewKademliaID("2111111400000000000000000000000000000000"), "localhost:8005")

	network2.SendFindContactMessage(&target)

	time.Sleep(5 * time.Second)

	for _, bucket := range network2.RoutingTable.buckets {
		for elt := bucket.list.Front(); elt != nil; elt = elt.Next() {
			contact := elt.Value.(Contact)
			fmt.Println(contact)
		}
	}
}

func TestPing(t *testing.T) {
	contact1 := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000")
	contact2 := NewContact(NewKademliaID("EFFFFFFF00000000000000000000000000000000"), "localhost:8001")
	kademlia1 := Kademlia{
		Network: &Network{
			RoutingTable: NewRoutingTable(contact1),
			Self:         &contact1,
		},
	}

	kademlia1.Network.RoutingTable.AddContact(NewContact(NewKademliaID("1111111200000000000000000000000000000000"), "localhost:8002"))
	kademlia1.Network.RoutingTable.AddContact(contact2)
	kademlia1.Network.RoutingTable.AddContact(NewContact(NewKademliaID("1111111300000000000000000000000000000000"), "localhost:8003"))
	kademlia1.Network.RoutingTable.AddContact(NewContact(NewKademliaID("1111111400000000000000000000000000000000"), "localhost:8004"))

	kademlia2 := Kademlia{
		Network: &Network{
			RoutingTable: NewRoutingTable(contact2),
			Self:         &contact2,
		},
	}

	kademlia2.Network.RoutingTable.AddContact(contact1)

	network1 := kademlia1.Network

	go Listen("localhost", "8000", &kademlia1)
	go Listen("localhost", "8001", &kademlia2)
	time.Sleep(1 * time.Second)

	network1.SendPingMessage(&contact2)

	time.Sleep(3 * time.Second)
}
