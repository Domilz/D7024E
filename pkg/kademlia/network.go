package kademlia

import (
	"fmt"
	"net"
	"strconv"
)

type Network struct {
	kademlia Kademlia
}

func Listen(ip string, port int) {
	fmt.Println("Started")
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

	go handleMessage(c)

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

	// TODO
}

func (network *Network) SendPingMessage(contact *Contact) {
	// TODO

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

func handleMessage(c chan []byte) {
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
