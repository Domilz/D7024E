package kademlia

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

type Network struct {
	Kademlia *Kademlia
}

type RPC struct {
	Topic    string      `json:"Topic"`
	Contact  Contact     `json:"Contact"`
	TargetID *KademliaID `json:"TargetID"`
}

func NewNetwork(kademlia *Kademlia) *Network {
	return &Network{Kademlia: kademlia}
}

func Listen(ip string, port int, kademliaNode *Network) error {
	address := ip + ":" + strconv.Itoa(port)

	conn, err := UDPServer(address)

	if err != nil {
		return err
	}

	defer conn.Close()

	buffer := make([]byte, 1024)
	c := make(chan []byte)
	remoteChan := make(chan *net.UDPAddr)

	go kademliaNode.handleMessages(c, remoteChan)

	for {
		// Read data from the UDP connection
		n, rAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("error reading from UDP connection:", err)
			continue
		}

		c <- buffer[:n]
		remoteChan <- rAddr

		// Print the received message
		fmt.Printf("Received %d bytes from %s: %s\n", n, rAddr, string(buffer[:n]))
	}

}

func (network *Network) SendPingMessage(contact *Contact) {
	conn := UDPConnection(contact.Address)
	fmt.Println("Sending ping to address,", contact.Address)
	defer conn.Close()

	ch := make(chan []byte)
	go func() {
		for {
			buffer := make([]byte, 1024)
			n, _, err := conn.ReadFromUDP(buffer)
			fmt.Println("Our conn:", conn.LocalAddr())
			if err != nil {
				fmt.Println("Error reading from udp server: ", err)
				return
			}
			ch <- buffer[:n]
		}
	}()

	rpc := RPC{
		Topic:   "ping",
		Contact: *network.Kademlia.Self,
	}

	message, err := json.Marshal(rpc)
	if err != nil {
		fmt.Println("Marshal ping error", err)
		return
	}

	_, err = conn.Write(message)
	if err != nil {
		fmt.Println("Error sending UDP message:", err)
		return
	}

	fmt.Println("Message sent to UDP server:", rpc)

	timeout := 4 * time.Second

	select {
	case data := <-ch:
		network.Kademlia.RoutingTable.UpdatedRecency(*contact)
		fmt.Println("Recieved pong: ", data)

	case <-time.After(timeout):
		fmt.Println("Timeout")
		network.Kademlia.RoutingTable.RemoveContact(*contact)
		network.Kademlia.RoutingTable.AddContact(rpc.Contact)
	}
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	closestNeighbour := network.Kademlia.RoutingTable.FindClosestContacts(contact.ID, 1)[0]
	conn := UDPConnection(closestNeighbour.Address)

	defer conn.Close()

	rpc := RPC{
		Topic:    "find_node",
		Contact:  *network.Kademlia.Self,
		TargetID: contact.ID,
	}

	message, err := json.Marshal(rpc)
	if err != nil {
		fmt.Println("Json error", err)
	}

	_, err = conn.Write(message)
	if err != nil {
		fmt.Println("Error sending UDP message:", err)
		return
	}

	fmt.Println("Message sent to UDP server:", string(message))

}

func UDPConnection(address string) *net.UDPConn {
	serverAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		fmt.Println(fmt.Errorf("error resolving UDPaddr for address: %s, error: %w", address, err))
	}

	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		fmt.Println(fmt.Errorf("error UDP dial to server: %s, error: %w", serverAddr, err))
	}

	return conn
}

func UDPServer(address string) (*net.UDPConn, error) {
	serverAddr, err := net.ResolveUDPAddr("udp4", address)

	if err != nil {
		return nil, fmt.Errorf("Error resolving UDPAddr: %w", err)
	}
	fmt.Println("Server address:", serverAddr)

	conn, err := net.ListenUDP("udp", serverAddr)
	if err != nil {
		return nil, fmt.Errorf("error UDP server listening: %w", err)
	}

	return conn, nil
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}

func (network *Network) handleMessages(c chan []byte, remoteChan chan *net.UDPAddr) {
	for {
		select {
		// Kollar ifall vi fått meddelande på kanalen
		case msg := <-c:
			var sender RPC
			rAddr := <-remoteChan

			err := json.Unmarshal(msg, &sender)
			if err != nil {
				fmt.Println("Unmarshal error: ", err)
			}
			fmt.Println("Topic:", sender.Topic)
			switch sender.Topic {
			case "ping":
				fmt.Println("Recieved a ping message!")
				fmt.Println("senderContact: ", sender.Contact)
				HandlePing(rAddr)
			case "find_node":
				go network.findNode(sender)

			case "find_node_response":
				go network.findNodeResponse(sender)
			default:
				fmt.Println("Other message:", string(msg))
			}
		}
	}
}

func (network *Network) findNode(rpc RPC) {
	closestNeigbourList := network.Kademlia.RoutingTable.FindClosestContacts(rpc.TargetID, 2)

	conn := UDPConnection(rpc.Contact.Address)

	defer conn.Close()

	for _, neigbour := range closestNeigbourList {
		newRPC := RPC{
			Topic:   "find_node_response",
			Contact: neigbour,
		}

		msg, err := json.Marshal(newRPC)
		if err != nil {
			fmt.Println("Marshal error: ", err)
		}
		conn.Write(msg)
		time.Sleep(1 * time.Second)
	}
}

func (network *Network) findNodeResponse(rpc RPC) {
	bucketIndex := network.Kademlia.RoutingTable.getBucketIndex(rpc.Contact.ID)
	bucket := network.Kademlia.RoutingTable.buckets[bucketIndex]
	var contact Contact
	if bucket.Len() >= bucketSize {
		contact = bucket.list.Back().Value.(Contact)
		network.SendPingMessage(&contact)
		return
	}
	network.Kademlia.RoutingTable.AddContact(rpc.Contact)
}

func GetLocalIP() string {
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println(fmt.Errorf("error hostname lookup: %w", err))
	}

	localIP, err := net.LookupIP(hostname)
	if err != nil {
		fmt.Println(fmt.Errorf("error IP lookup: %w", err))
	}

	return string(localIP[0])
}

func HandlePing(address *net.UDPAddr) {
	fmt.Println("Handle ping")
	conn, err := net.DialUDP("udp", nil, address)
	defer conn.Close()
	if err != nil {
		fmt.Println(fmt.Errorf("error in connection when handling ping: %w", err))
	}

	fmt.Println("Sending pong to ", address)
	conn.Write([]byte("I live"))
}
