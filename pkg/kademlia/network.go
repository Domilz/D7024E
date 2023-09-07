package kademlia

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

type Network struct {
	Kademlia *Kademlia
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

	// c := make(chan []byte)

	for {
		buffer := make([]byte, 1024)
		n, rAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("error reading from UDP connection: ", err)
		}

		response := kademliaNode.handleMessages(buffer[:n])

		go kademliaNode.WriteResponse(response, rAddr, conn)

		fmt.Printf("Received %d bytes from %s: %s\n", n, rAddr, string(buffer[:n]))

	}

}

func (network *Network) sendMessage(address string, msg []byte) ([]byte, error) {
	conn := UDPConnect(address)

	defer conn.Close()

	ch := make(chan []byte)

	go func() {
		buffer := make([]byte, 1024)
		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println(fmt.Errorf("Dont worry only error reading from UDP in sendMessage: %w ", err))
			return
		}
		ch <- buffer[:n]
	}()

	_, err := conn.Write(msg)
	if err != nil {
		return nil, fmt.Errorf("Error sending UDP message: %w", err)
	}

	fmt.Println("Message sent to UDP server:", string(msg))

	timeout := 4 * time.Second

	select {
	case data := <-ch:
		return data, nil

	case <-time.After(timeout):
		return nil, errors.New("response timed out")
	}
}

func (network *Network) handleMessages(content []byte) RPC {
	var rpc RPC
	err := json.Unmarshal(content, &rpc)
	if err != nil {
		fmt.Println("Unmarshal handleMessage error", err)
	}

	switch rpc.Topic {
	case "ping":
		return network.CreateRPC("pong", *network.Kademlia.Self, nil, nil)
	case "find_node":
		return network.findNode(rpc)
	default:
		return RPC{}
	}
}

func (network *Network) WriteResponse(response RPC, rAddr *net.UDPAddr, conn *net.UDPConn) {
	byteResponse, err := json.Marshal(response)
	if err != nil {
		fmt.Println("Json marshal WriteResponse error", err)
	}
	_, err = conn.WriteTo(byteResponse, rAddr)
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func (network *Network) SendPingMessage(contact *Contact) bool {
	fmt.Println("Sending ping to address,", contact.Address)

	rpc := network.CreateRPC("ping", *contact, nil, nil)

	message, err := json.Marshal(rpc)
	if err != nil {
		fmt.Println("Marshal ping error", err)
	}

	response, err := network.sendMessage(contact.Address, message)
	if err != nil {
		fmt.Println("Pong failed:", err)
		return false

	} else {
		fmt.Println("Recieved pong:", response)
		return true
	}
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	closestNeighbour := network.Kademlia.RoutingTable.FindClosestContacts(contact.ID, 1)[0]

	rpc := network.CreateRPC("find_node", *network.Kademlia.Self, contact.ID, nil)

	message, err := json.Marshal(rpc)
	if err != nil {
		fmt.Println("Json error", err)
	}

	response, err := network.sendMessage(closestNeighbour.Address, message)
	fmt.Println("Message sent to UDP server:", string(message))

	if err != nil {
		fmt.Println("error: Node response timed out", err)
	} else {
		var responseRPC RPC
		err := json.Unmarshal(response, &responseRPC)
		if err != nil {
			fmt.Println("error unmarshal of node response:", err)
		}
		network.findNodeResponse(responseRPC)
	}

}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
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

func UDPConnect(address string) *net.UDPConn {
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
