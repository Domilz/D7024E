package kademlia

import "fmt"

type RPC struct {
	Topic       string      `json:"Topic"`
	Contact     Contact     `json:"Contact"`
	TargetID    *KademliaID `json:"TargetID"`
	ContactList []Contact   `json:"ContactList"`
}

func (network *Network) CreateRPC(topic string, contact Contact, targetID *KademliaID, contactList []Contact) RPC {
	return RPC{
		Topic:       topic,
		Contact:     contact,
		TargetID:    targetID,
		ContactList: contactList,
	}
}

func (network *Network) findNode(rpc RPC) RPC {
	closestNeigbourList := network.Kademlia.RoutingTable.FindClosestContacts(rpc.TargetID, 2)

	newRPC := network.CreateRPC("find_node_response", *network.Kademlia.Self, nil, closestNeigbourList)
	return newRPC
}

func (network *Network) findNodeResponse(rpc RPC) {
	for _, contact := range rpc.ContactList {
		bucketIndex := network.Kademlia.RoutingTable.getBucketIndex(contact.ID)
		bucket := network.Kademlia.RoutingTable.buckets[bucketIndex]

		var pingContact Contact
		if bucket.Len() >= bucketSize {
			pingContact = bucket.list.Back().Value.(Contact)
			flag := network.SendPingMessage(&pingContact)
			if flag {
				network.Kademlia.RoutingTable.UpdatedRecency(pingContact)
			} else {
				network.Kademlia.RoutingTable.RemoveContact(pingContact)
				fmt.Println("Remove contact done")
				network.Kademlia.RoutingTable.AddContact(contact)
				fmt.Println("Add contact done")
			}
		} else {
			network.Kademlia.RoutingTable.AddContact(contact)
		}
	}
}
