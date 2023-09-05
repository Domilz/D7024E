package kademlia

import (
	"fmt"
	"net"
	"os"
	"strconv"
)

type Network struct {
	Kademlia *Kademlia
}

func NewNetwork(kademlia *Kademlia) *Network {
	return &Network{Kademlia: kademlia}
}

func Listen(ip string, port int) {
	hname, err := os.Hostname()
	if err != nil {
		fmt.Println("Error", err)
	}
	localIP, err := net.LookupIP(hname)
	if err != nil {
		fmt.Println("Error", err)
	}
	fmt.Println("IP", localIP[0])

	p := ":" + strconv.Itoa(port)

	serverAddr, err := net.ResolveUDPAddr("udp4", p)
	if err != nil {
		fmt.Println(fmt.Errorf("err", err))
	}

	conn, err := net.ListenUDP("udp", serverAddr)
	if err != nil {
		fmt.Println(fmt.Errorf("err", err))
	}

	defer conn.Close()

	buffer := make([]byte, 1024)
	c := make(chan []byte)

	go handleMessages(c)

	for {
		// Read data from the UDP connection
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading from UDP connection:", err)
			continue
		}

		c <- buffer[:n]

		// Print the received message
		fmt.Printf("Received %d bytes from %s: %s\n", n, addr, string(buffer[:n]))
	}

}

func (network *Network) SendPingMessage() {
	// ips, err := net.LookupIP("kademlianodes")
	// if err != nil {
	// 	fmt.Println("Error resolving IP address:", err)
	// 	return
	// }

	// for _, ip := range ips {
	// 	fmt.Println(ip)
	// }

	bootstrapNodeHostname := os.Getenv("BOOTSTRAP_NODE_IP")
	ips, err := net.LookupIP(bootstrapNodeHostname)
	if err != nil {
		fmt.Println("Error", err)
	}
	ip := ips[0].String() + ":" + os.Getenv("NODE_PORT")

	serverAddr, err := net.ResolveUDPAddr("udp", ip)
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

	message := []byte("ping")

	_, err = conn.Write(message)
	if err != nil {
		fmt.Println("Error sending UDP message:", err)
		return
	}

	fmt.Println("Message sent to UDP server:", string(message))
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}

func handleMessages(c chan []byte) {
	for {
		select {
		// Kollar ifall vi fått meddelande på kanalen
		case msg := <-c:
			switch string(msg) {
			case "ping":
				fmt.Println("Recieved a ping message!")
			}
		}
	}
}
