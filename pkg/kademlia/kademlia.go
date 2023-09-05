package kademlia

import (
	"fmt"
	"net"
	"os"
)

type Kademlia struct {
	RoutingTable *RoutingTable
	Self         *Contact
}

const (
	DefaultBootstrapInput = "FFFFFFFF00000000000000000000000000000000"
)

func NewKademlia() *Kademlia {

	container := os.Getenv("ISBOOTSTRAP")
	var contact Contact
	switch container {
	case "1":
		bootstrapNodeHostname := os.Getenv("BOOTSTRAP_NODE_IP")
		ips, err := net.LookupIP(bootstrapNodeHostname)
		if err != nil {
			fmt.Println("Error", err)
		}
		ip := ips[0].String() + ":" + os.Getenv("NODE_PORT")
		contact = NewContact(NewKademliaID(DefaultBootstrapInput), ip)
	default:
		contact = NewContact(NewRandomKademliaID(), "localhost:8080")
	}
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
