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

func (kademlia Kademlia) LookupContact(target Contact) []Contact {
	closeToTarget := kademlia.Network.RoutingTable.FindClosestContacts(target.ID, K)
	var closestContacts []Nodes
	for i, node := range closeToTarget {
		fmt.Println("Node:", node)
		newNode := Nodes{contact: &closeToTarget[i], visited: false}
		closestContacts = append(closestContacts, newNode)
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	kademlia.getClosestFromLookup(ctx, cancel, &closestContacts, target)

	fmt.Println("")
	for _, ele := range closestContacts {
		fmt.Println(ele.contact, ele.visited)
	}

	return nil
}

func (kademlia Kademlia) getClosestFromLookup(ctx context.Context, cancel context.CancelFunc, closestContacts *[]Nodes, target Contact) []Nodes {
	responseChannel := make(chan []Contact)
	doneCh := make(chan bool)
	go kademlia.sendFindNode(responseChannel, doneCh, closestContacts, target, Alpha)

	newSends := 0
	for {
		select {
		case contactList := <-responseChannel:
			writeLock.Lock()
			newElement := kademlia.UpdateContacts(closestContacts, contactList, target)
			writeLock.Unlock()
			if newElement {
				newSends++
				go kademlia.getClosestFromLookup(ctx, cancel, closestContacts, target)
			}
		case <-doneCh:
			if newSends == 0 {
				writeLock.Lock()
				kademlia.finalLookup(closestContacts, target)
				writeLock.Unlock()
				cancel()
			}

		case <-ctx.Done():
			return nil
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
			_ = kademlia.UpdateContacts(closestContacts, contactList, target)
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

func (kademlia Kademlia) UpdateContacts(closestContacts *[]Nodes, potentialContacts []Contact, target Contact) bool {
	var newList []Nodes
	i := 0
	j := 0
	addedNewElement := false

	for len(newList) < K && (i < len(*closestContacts) || j < len(potentialContacts)) {
		if j == len(potentialContacts) || !potentialContacts[j].ID.CalcDistance(target.ID).Less((*closestContacts)[i].contact.distance) {
			newList = append(newList, (*closestContacts)[i])
			i++
		} else {
			potentialContacts[j].distance = potentialContacts[j].ID.CalcDistance(target.ID)
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

// func checkVisited(closesContacts []Nodes) bool {
// 	for _, c := range closesContacts {
// 		if !c.visited {
// 			return false
// 		}
// 	}
// 	return true
// }

// 	time.Sleep(1 * time.Second)

// 	for {
// 		fmt.Println("waiting for selext")
// 		select {
// 		case contactList := <-contactResponseChannel:
// 			fmt.Println("contactList: ", contactList)
// 			fmt.Println("Closest contacts: ", testClosest)
// 			// Update closestContactsSoFar
// 			testClosest, addedNewElement := kademlia.UpdateContacts(testClosest, contactList, target)
// 			fmt.Println("After loop closest contacts: ", testClosest)

// 			if addedNewElement {
// 				go kademlia.resendFindNode(contactResponseChannel, &testClosest, Alpha, target)
// 			} else {
// 				if checkVisited(closestContacts) {
// 					return closestContacts
// 				}
// 				go kademlia.resendFindNode(contactResponseChannel, &testClosest, K, target)
// 			}
// 		}
// 	}
// }

// func checkVisited(closesContacts []NodesVisited) bool {
// 	for _, c := range closesContacts {
// 		if !c.visited {
// 			return false
// 		}
// 	}
// 	return true
// }

// func (kademlia Kademlia) resendFindNode(responseChannel chan []Contact, closestContacts *[]NodesVisited, times int, target Contact) {
// 	fmt.Println("I am in")
// 	i := 0
// 	j := 0
// 	for i < times && j < len(*closestContacts) {
// 		if (*closestContacts)[j].visited {
// 			contactList := kademlia.Network.SendFindContactMessage(&(*closestContacts)[j].contact, target.ID)
// 			(*closestContacts)[j].visited = true
// 			responseChannel <- contactList
// 			i++
// 		}
// 		j++
// 	}
// 	fmt.Println("I am done")
// }

// func (kademlia Kademlia) UpdateContacts(closestContacts []NodesVisited, potentialContacts []Contact, target Contact) ([]NodesVisited, bool) {
// 	var newList []NodesVisited
// 	i := 0
// 	j := 0
// 	addedNewElement := false

// 	for len(newList) < K && (i < len(closestContacts) || j < len(potentialContacts)) {
// 		// if j == len(potentialContacts) || closestContacts[i].contact.distance.Less(potentialContacts[j].ID.CalcDistance(target.ID)) {
// 		if j == len(potentialContacts) || !potentialContacts[j].ID.CalcDistance(target.ID).Less(closestContacts[i].contact.distance) {
// 			newList = append(newList, closestContacts[i])
// 			i++
// 		} else {
// 			potentialContacts[j].distance = potentialContacts[j].ID.CalcDistance(target.ID)
// 			newList = append(newList, NodesVisited{contact: (potentialContacts)[j], visited: false})
// 			addedNewElement = true
// 			j++
// 		}
// 	}
// 	return newList, addedNewElement
// }

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}

// func (kademlia Kademlia) LookupContact(target Contact) []Contact {
// 	neighbours := kademlia.Network.RoutingTable.FindClosestContacts(target.ID, K)
// 	var closestContacts []NodesVisited
// 	for _, node := range neighbours {
// 		fmt.Println("Node:", node)
// 		closestContacts = append(closestContacts, NodesVisited{contact: node, visited: false})
// 	}
// 	nodes := kademlia.getClosestFromLookup(closestContacts, target)

// 	var returnContacts []Contact
// 	for i, node := range nodes {
// 		fmt.Println(i, node.visited)
// 		returnContacts = append(returnContacts, node.contact)
// 	}
// 	return returnContacts
// }

// func (kademlia Kademlia) getClosestFromLookup(closestContacts []NodesVisited, target Contact) []NodesVisited {
// 	testClosest := closestContacts
// 	contactResponseChannel := make(chan []Contact)
// 	go func() {
// 		for i := 0; i < Alpha; i++ {
// 			contact := testClosest[i]
// 			contactList := kademlia.Network.SendFindContactMessage(&contact.contact, target.ID)
// 			fmt.Printf("Contact %v got Contactlist: %v\n", contact, contactList)
// 			contact.visited = true
// 			contactResponseChannel <- contactList
// 		}
// 	}()
