package kademlia

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFindNode(t *testing.T) {

	contact1 := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000")
	contact2 := NewContact(NewKademliaID("EFFFFFFF00000000000000000000000000000000"), "localhost:8001")
	kademlia1 := Kademlia{
		Network: &Network{
			RoutingTable: NewRoutingTable(contact1),
			Self:         &contact1,
		},
	}

	addContact1 := NewContact(NewKademliaID("0000000000000000000000000000000000000001"), "localhost:8003")
	addContact2 := NewContact(NewKademliaID("E000000000000000000000000000000000000003"), "localhost:8004")
	addContact3 := NewContact(NewKademliaID("0000000000000000000000000000000000000002"), "localhost:8002")
	addContact4 := NewContact(NewKademliaID("0000000000000000000000000000000000000004"), "localhost:8002")
	kademlia1.Network.RoutingTable.AddContact(addContact3)
	kademlia1.Network.RoutingTable.AddContact(addContact1)
	kademlia1.Network.RoutingTable.AddContact(addContact2)
	kademlia1.Network.RoutingTable.AddContact(addContact4)

	kademlia2 := Kademlia{
		Network: &Network{
			RoutingTable: NewRoutingTable(contact2),
			Self:         &contact2,
		},
	}

	kademlia2.Network.RoutingTable.AddContact(contact1)

	network2 := kademlia2.Network

	wantedList := []Contact{addContact1, addContact3, addContact4}

	go Listen("localhost", "8000", &kademlia1)
	go Listen("localhost", "8001", &kademlia2)
	time.Sleep(1 * time.Second)

	contactList := network2.SendFindContactMessage(&contact1, addContact1.ID)

	time.Sleep(2 * time.Second)

	assert.Equal(t, wantedList, contactList)

	fmt.Println(contactList)

	// for _, bucket := range network2.RoutingTable.buckets {
	// 	for elt := bucket.list.Front(); elt != nil; elt = elt.Next() {
	// 		contact := elt.Value.(Contact)
	// 		fmt.Println(contact)
	// 	}
	// }
}

func TestStore(t *testing.T) {
	contact1 := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000")
	contact2 := NewContact(NewKademliaID("EFFFFFFF00000000000000000000000000000000"), "localhost:8001")
	kademlia1 := Kademlia{
		Network: &Network{
			RoutingTable: NewRoutingTable(contact1),
			Self:         &contact1,
			Objects:      make(map[KademliaID]string),
		},
	}

	kademlia2 := Kademlia{
		Network: &Network{
			RoutingTable: NewRoutingTable(contact2),
			Self:         &contact2,
			Objects:      make(map[KademliaID]string),
		},
	}

	network1 := kademlia1.Network

	go Listen("localhost", "8000", &kademlia1)
	go Listen("localhost", "8001", &kademlia2)
	time.Sleep(1 * time.Second)

	hash := network1.SendStoreMessage(contact2, []byte("hehe"))
	time.Sleep(1 * time.Second)
	value := kademlia2.Network.Objects[*hash]
	fmt.Println(value)
	assert.Equal(t, "hehe", value)

}

func TestListen(t *testing.T) {

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

	pingWork := network1.SendPingMessage(&contact2)
	assert.True(t, pingWork)

	time.Sleep(3 * time.Second)
}
