package main

import (
	"fmt"
	"net"
)

func main() {
	// x := kademlia.NewKademlia()
	fmt.Println("Started")
	server, err := net.ResolveUDPAddr("udp4", ":8000")
	if err != nil {
		fmt.Println(fmt.Errorf("err", err))
	}
	conn, err := net.ListenUDP("udp", server)

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
