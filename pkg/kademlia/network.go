package kademlia

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
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

func Listen(ip string, port int) error {

	kademliaNode := NewNetwork(NewKademlia())

	address := ip + ":" + strconv.Itoa(port)

	conn, err := UDPServer(address)

	if err != nil {
		return err
	}

	defer conn.Close()

	buffer := make([]byte, 1024)
	c := make(chan []byte)

	go kademliaNode.handleMessages(c)

	for {
		// Read data from the UDP connection
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("error reading from UDP connection:", err)
			continue
		}

		c <- buffer[:n]

		// Print the received message
		fmt.Printf("Received %d bytes from %s: %s\n", n, addr, string(buffer[:n]))
	}

}

func (network *Network) SendPingMessage(contact *Contact) {
	conn := UDPConnection(contact.Address)

	defer conn.Close()
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

func UDPConnection(address string) net.Conn {
	serverAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		fmt.Println(fmt.Errorf("Error resolving UDPaddr for address: %s, error: %w\n", address, err))
	}

	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		fmt.Println(fmt.Errorf("Error UDP dial to server: %s, error: %w", serverAddr, err))
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

func (network *Network) handleMessages(c chan []byte) {
	var sender RPC
	for {
		select {
		// Kollar ifall vi fått meddelande på kanalen
		case msg := <-c:
			err := json.Unmarshal(msg, &sender)
			if err != nil {
				fmt.Println("Unmarshal error: ", err)
			}
			switch sender.Topic {
			case "ping":
				fmt.Println("Recieved a ping message!")
				fmt.Println("senderContact: ", sender.Contact)
			case "find_node":
				network.findNode(sender)

			case "find_node_response":
				fmt.Println("response: ", sender)
			default:
				fmt.Println("Other message:", string(msg))
			}
		}
	}
}

func (network *Network) findNode(rpc RPC) {
	closestNeigbourList := network.Kademlia.RoutingTable.FindClosestContacts(rpc.TargetID, 1)
	serverAddr, err := net.ResolveUDPAddr("udp", rpc.Contact.Address)
	if err != nil {
		fmt.Println("Error resolving server address:", err)
		return
	}

	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		fmt.Println("Error creating UDP connection:", err)
		return
	}

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

	}

}
