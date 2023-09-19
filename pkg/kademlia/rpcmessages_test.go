package kademlia

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStoreRPC(t *testing.T) {
	contact1 := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000")
	kademlia1 := Kademlia{
		Network: &Network{
			RoutingTable: NewRoutingTable(contact1),
			Self:         &contact1,
			Objects:      make(map[KademliaID]string),
		},
	}

	key := kademlia1.Network.AddValueToDic("hej")
	fmt.Println(key)
	value := kademlia1.Network.Objects[*key]

	assert.Equal(t, "hej", value)
}
