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
	key := network.AddValueToDic(rpc.Value)
	newRPC := network.CreateRPC("store_response", *network.Self, key, nil, "")
	return newRPC
}

func (network *Network) findNode(rpc RPC) RPC {
	closestNeigbourList := network.RoutingTable.FindClosestContacts(rpc.TargetID, K)

	newRPC := network.CreateRPC("find_node_response", *network.Self, nil, closestNeigbourList, "")
	return newRPC
}

func (network *Network) findValue(rpc RPC) RPC {
	byteRep := []byte(rpc.Value)
	var key *KademliaID
	for i := 0; i < IDLength; i++ {
		(*key)[i] = byteRep[i]
	}
	val, ok := network.Objects[*rpc.TargetID]
	if ok {
		return network.CreateRPC("find_value_response", *network.Self, rpc.TargetID, nil, val)
	}
	contactList := network.RoutingTable.FindClosestContacts(key, K)

	return network.CreateRPC("find_value_response", *network.Self, nil, contactList, rpc.Value)
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
