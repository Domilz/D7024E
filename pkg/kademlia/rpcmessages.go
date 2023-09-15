package kademlia

type RPC struct {
	Topic       string      `json:"Topic"`
	Contact     Contact     `json:"Contact"`
	TargetID    *KademliaID `json:"TargetID"`
	ContactList []Contact   `json:"ContactList"`
	Value       string      `json:"Value"`
}

func (network *Network) CreateRPC(topic string, contact Contact, targetID *KademliaID, contactList []Contact, value string) RPC {
	return RPC{
		Topic:       topic,
		Contact:     contact,
		TargetID:    targetID,
		ContactList: contactList,
		Value:       value,
	}
}

func (network *Network) StoreValue(rpc RPC) RPC {
	network.AddValueToDic(rpc.Value)
	newRPC := network.CreateRPC("store_response", *network.Self, nil, nil, "")
	return newRPC
}

func (network *Network) findNode(rpc RPC) RPC {
	closestNeigbourList := network.RoutingTable.FindClosestContacts(rpc.TargetID, bucketSize)

	newRPC := network.CreateRPC("find_node_response", *network.Self, nil, closestNeigbourList, "")
	return newRPC
}

func (network *Network) findValue(rpc RPC) RPC {
	network.RoutingTable.FindClosestContacts(rpc.TargetID, bucketSize)
}

func (network *Network) updateRoutingTable(rpc RPC) {
	for _, contact := range rpc.ContactList {
		bucketIndex := network.RoutingTable.getBucketIndex(contact.ID)
		bucket := network.RoutingTable.buckets[bucketIndex]

		var pingContact Contact
		if bucket.Len() >= bucketSize {
			pingContact = bucket.list.Back().Value.(Contact)
			flag := network.SendPingMessage(&pingContact)
			if flag {
				network.RoutingTable.UpdatedRecency(pingContact)
			} else {
				network.RoutingTable.RemoveContact(pingContact)
				network.RoutingTable.AddContact(contact)
			}
		} else {
			network.RoutingTable.AddContact(contact)
		}
	}
}
