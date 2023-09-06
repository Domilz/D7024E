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
	var contact Contact
	var routingTable *RoutingTable

	container := os.Getenv("ISBOOTSTRAP")
	bootstrapNodeHostname := os.Getenv("BOOTSTRAP_NODE_IP")
	ips, err := net.LookupIP(bootstrapNodeHostname)
	if err != nil {
		fmt.Println("error for boostrap node IP lookup:", err)
	}
	ip := ips[0].String() + ":" + os.Getenv("NODE_PORT")
	bootstrap_contact := NewContact(NewKademliaID(DefaultBootstrapInput), ip)

	switch container {
	case "1":
		contact = bootstrap_contact
		routingTable = NewRoutingTable(contact)
	default:
		hostname, err := os.Hostname()
		if err != nil {
			fmt.Println("error hostname lookup for kademlia node:", err)
		}
		localIP, err := net.LookupIP(hostname)
		if err != nil {
			fmt.Println("error IP lookup for kademlia node:", err)
		}

		contact = NewContact(NewRandomKademliaID(), localIP[0].String()+":8080")
		routingTable = NewRoutingTable(contact)
		routingTable.AddContact(bootstrap_contact)
		closestContacts := routingTable.FindClosestContacts(contact.ID, 1)
		fmt.Println("closes contact: ", closestContacts)
	}
	return &Kademlia{
		Self:         &contact,
		RoutingTable: routingTable,
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
