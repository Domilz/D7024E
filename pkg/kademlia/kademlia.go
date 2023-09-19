package kademlia

import (
	"context"
	"fmt"
	"net"
	"os"
	"sync"
)

type Kademlia struct {
	Network *Network
	Objects map[KademliaID]string
}

type Nodes struct {
	contact *Contact
	visited bool
}

const (
	DefaultBootstrapInput = "FFFFFFFF00000000000000000000000000000000"
	Alpha                 = 3
	K                     = 20
)

var (
	writeLock sync.Mutex
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

func (kademlia *Kademlia) JoinNetwork() {
	_ = kademlia.LookupContact(*kademlia.Network.Self)
}

func (kademlia Kademlia) LookupContact(target Contact) []Contact {
	closeToTarget := kademlia.Network.RoutingTable.FindClosestContacts(target.ID, K)
	var closestContacts []Nodes
	for i, _ := range closeToTarget {
		newNode := Nodes{contact: &closeToTarget[i], visited: false}
		closestContacts = append(closestContacts, newNode)
	}

	mainCh := make(chan bool)
	go kademlia.getClosestFromLookup(mainCh, &closestContacts, target)
	<-mainCh
	writeLock.Lock()
	kademlia.finalLookup(&closestContacts, target)
	writeLock.Unlock()

	var returnList []Contact
	for _, ele := range closestContacts {
		returnList = append(returnList, *ele.contact)
	}

	return returnList
}

func (kademlia Kademlia) getClosestFromLookup(finishedCh chan bool, closestContacts *[]Nodes, target Contact) {
	responseChannel := make(chan []Contact)
	doneCh := make(chan bool)
	roundsCh := make(chan bool)
	go kademlia.sendFindNode(responseChannel, doneCh, closestContacts, target, Alpha)

	newRounds := 0
	for {
		select {
		case contactList := <-responseChannel:
			writeLock.Lock()
			newElement := kademlia.UpdateContacts(closestContacts, contactList, target.ID)
			writeLock.Unlock()
			if newElement {
				newRounds++
				go kademlia.getClosestFromLookup(roundsCh, closestContacts, target)
			}
		case <-doneCh:
			for i := 0; i < newRounds; i++ {
				<-roundsCh
			}
			finishedCh <- true
			return
		}
	}
}

func (kademlia Kademlia) finalLookup(closestContacts *[]Nodes, target Contact) {
	// contact := NewContact(NewKademliaID("1111111600000000000000000000000000000000"), "localhost:8005")
	// *closestContacts = append(*closestContacts, Nodes{contact: &contact, visited: false})
	i := 0
	for i < len(*closestContacts) {
		if !(*closestContacts)[i].visited {
			(*closestContacts)[i].visited = true
			contactList := kademlia.Network.SendFindContactMessage((*closestContacts)[i].contact, target.ID)
			_ = kademlia.UpdateContacts(closestContacts, contactList, target.ID)
		}
		i++
	}
}

func (kademlia Kademlia) sendFindNode(responseChannel chan []Contact, doneCh chan bool, closestContacts *[]Nodes, target Contact, times int) {
	i := 0
	j := 0
	for i < times && j < len(*closestContacts) {
		if !(*closestContacts)[j].visited {
			(*closestContacts)[j].visited = true
			contactList := kademlia.Network.SendFindContactMessage((*closestContacts)[j].contact, target.ID)
			responseChannel <- contactList
			i++
		}
		j++
	}
	doneCh <- true
}

func (kademlia Kademlia) UpdateContacts(closestContacts *[]Nodes, potentialContacts []Contact, target *KademliaID) bool {
	var newList []Nodes
	i := 0
	j := 0
	addedNewElement := false

	for len(newList) < K && (i < len(*closestContacts) || j < len(potentialContacts)) {
		if j == len(potentialContacts) || !potentialContacts[j].ID.CalcDistance(target).Less((*closestContacts)[i].contact.distance) {
			newList = append(newList, (*closestContacts)[i])
			i++
		} else {
			potentialContacts[j].distance = potentialContacts[j].ID.CalcDistance(target)
			if !Find(&newList, potentialContacts[j]) {
				newList = append(newList, Nodes{contact: &(potentialContacts)[j], visited: false})
			}
			addedNewElement = true
			j++
		}
	}

	*closestContacts = newList

	return addedNewElement
}

func Find(list *[]Nodes, contact Contact) bool {
	for _, ele := range *list {
		if ele.contact.ID.Equals(contact.ID) {
			return true
		}
	}
	return false
}

func (kademlia *Kademlia) LookupData(hash string) string {
	byteRepresentation := []byte(hash)
	var target KademliaID
	for i := 0; i < IDLength; i++ {
		target[i] = byteRepresentation[i]
	}
	closeToTarget := kademlia.Network.RoutingTable.FindClosestContacts(&target, K)
	var closestContacts []Nodes
	for i, _ := range closeToTarget {
		newNode := Nodes{contact: &closeToTarget[i], visited: false}
		closestContacts = append(closestContacts, newNode)
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	kademlia.getClosestFromLookupData(ctx, cancel, &closestContacts, &target, hash)

	for _, ele := range closestContacts {
		fmt.Println(ele.contact, ele.visited)
	}

	return ""
}

func (kademlia Kademlia) getClosestFromLookupData(ctx context.Context, cancel context.CancelFunc, closestContacts *[]Nodes, target *KademliaID, data string) string {
	responseChannel := make(chan []Contact)
	doneCh := make(chan string)
	go kademlia.sendFindData(responseChannel, doneCh, closestContacts, data, Alpha)

	newSends := 0
	for {
		select {
		case contactList := <-responseChannel:
			writeLock.Lock()
			newElement := kademlia.UpdateContacts(closestContacts, contactList, target)
			writeLock.Unlock()
			if newElement {
				newSends++
				go kademlia.getClosestFromLookupData(ctx, cancel, closestContacts, target, data)
			}
		case val := <-doneCh:
			if val != "not found" {
				cancel()
				return val
			}
			return "not here"

		case <-ctx.Done():
			return "found"
		}
	}
}

func (kademlia Kademlia) sendFindData(responseChannel chan []Contact, doneCh chan string, closestContacts *[]Nodes, data string, times int) {
	i := 0
	j := 0
	for i < times && j < len(*closestContacts) {
		if !(*closestContacts)[j].visited {
			(*closestContacts)[j].visited = true
			contactList, value := kademlia.Network.SendFindDataMessage((*closestContacts)[j].contact, data)
			if contactList == nil {
				doneCh <- value
			}
			responseChannel <- contactList
			i++
		}
		j++
	}
	doneCh <- "not found"
}

func (kademlia *Kademlia) Store(data []byte) (KademliaID, error) {
	key := NewKademliaID(string(data))
	c := NewContact(key, "")
	contacts := kademlia.LookupContact(c)
	successfully := false
	for _, contact := range contacts {
		_, err := kademlia.Network.SendStoreMessage(contact, data)
		if err == nil {
			successfully = true
		}
	}

	if successfully {
		return *key, nil
	}
	return *key, fmt.Errorf("uploading failed")
}
