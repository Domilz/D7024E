package kademlia

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"time"
)

type Network struct {
	RoutingTable *RoutingTable
	Self         *Contact
	Objects      map[KademliaID]string
}

func Listen(ip string, port string, kademliaNode *Kademlia) error {
	address := ip + ":" + port

	conn, err := UDPServer(address)

	if err != nil {
		return err
	}

	defer conn.Close()

	for {
		buffer := make([]byte, 1024)
		n, rAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("error reading from UDP connection: ", err)
		}

		response := kademliaNode.Network.handleMessages(buffer[:n])

		go kademliaNode.Network.WriteResponse(response, rAddr, conn)

		// fmt.Printf("Received %d bytes from %s: %s\n", n, rAddr, string(buffer[:n]))
		//fmt.Printf("Received %d bytes from %s: %s\n", n, rAddr, string(buffer[:n]))
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
		return nil, fmt.Errorf("error sending UDP message: %w", err)
	}

	// fmt.Println("Message sent to UDP server:", string(msg))
	//fmt.Println("Message sent to UDP server:", string(msg))

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

	fmt.Println("Topic: ", rpc.Topic)
	switch rpc.Topic {
	case "ping":
		fmt.Println("in ping")
		return network.CreateRPC("pong", *network.Self, nil, nil, "")
	case "find_node":
		return network.findNode(rpc)
	case "store_value":
		return network.StoreValue(rpc)
	case "find_data":
		return network.FindValue(rpc)
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

	rpc := network.CreateRPC("ping", *contact, nil, nil, "")

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

func (network *Network) SendFindContactMessage(contact *Contact, targetID *KademliaID) []Contact {

	rpc := network.CreateRPC("find_node", *network.Self, targetID, nil, "")

	message, err := json.Marshal(rpc)
	if err != nil {
		fmt.Println("Json error", err)
	}

	response, err := network.sendMessage(contact.Address, message)
	// fmt.Println("Message sent to UDP server:", contact.Address, "With message:", string(message))

	if err != nil {
		fmt.Println("error: Node response timed out", err)
	} else {
		var responseRPC RPC
		err := json.Unmarshal(response, &responseRPC)
		if err != nil {
			fmt.Println("error unmarshal of node response:", err)
		}
		network.updateRoutingTable(responseRPC)

		return responseRPC.ContactList
	}
	return nil
}

func (network *Network) SendFindDataMessage(contact *Contact, hash string) ([]Contact, string) {
	rpc := network.CreateRPC("find_data", *network.Self, nil, nil, hash)

	message, err := json.Marshal(rpc)
	if err != nil {
		fmt.Println("Json error", err)
	}

	response, err := network.sendMessage(contact.Address, message)
	//fmt.Println("Message sent to UDP server:", string(message))

	if err != nil {
		fmt.Println("error: Node response timed out", err)
	} else {
		var responseRPC RPC
		err := json.Unmarshal(response, &responseRPC)
		if err != nil {
			fmt.Println("error unmarshal of node response:", err)
		}

		if responseRPC.ContactList == nil {
			return nil, responseRPC.Value
		} else {
			network.updateRoutingTable(responseRPC)
			return responseRPC.ContactList, ""
		}

	}
	return nil, ""
}

func (network *Network) SendStoreMessage(contact Contact, data []byte) (*KademliaID, error) {
	rpc := network.CreateRPC("store_value", *network.Self, nil, nil, string(data))

	message, err := json.Marshal(rpc)
	if err != nil {
		fmt.Println("Json error", err)
	}

	response, err := network.sendMessage(contact.Address, message)
	//fmt.Println("Message sent to UDP server:", string(message))

	if err != nil {
		return nil, fmt.Errorf("Store failed: %v", err)

	} else {
		var responseRPC RPC
		err := json.Unmarshal(response, &responseRPC)
		if err != nil {
			fmt.Println("error unmarshal of node response:", err)
		}

		return responseRPC.TargetID, nil

	}

}

func (network *Network) AddValueToDic(value string) *KademliaID {
	var key KademliaID
	key = sha1.Sum([]byte(value))

	network.Objects[key] = value
	return &key
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
	serverAddr, err := net.ResolveUDPAddr("udp", address)

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
