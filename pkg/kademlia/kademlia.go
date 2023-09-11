package kademlia

import (
	"fmt"
	"net"
	"os"
)

type Kademlia struct {
	Network *Network
	Objects map[KademliaID]string
}

const (
	DefaultBootstrapInput = "FFFFFFFF00000000000000000000000000000000"
	Alpha                 = 3
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
		Network: &Network{
			Self:         &contact,
			RoutingTable: routingTable,
		},
	}
}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	neighbours := kademlia.Network.RoutingTable.FindClosestContacts(target.ID, Alpha)
	kademlia.getClosestFromLookup(neighbours, target)
}

func (kademlia *Kademlia) getClosestFromLookup(closestContactsSoFar []Contact, target *Contact) []Contact {
	contactResponseChanel := make(chan []Contact)
	for _, contact := range closestContactsSoFar {
		go func() {

			contactsList := kademlia.Network.SendFindContactMessage(&contact)
			contactResponseChanel <- contactsList
		}()
	}

	for {

	}
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
