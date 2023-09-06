package main

import (
	"fmt"

	"github.com/Domilz/D7024E/pkg/kademlia"
)

func main() {
	fmt.Println("Starting...")

	localIP := kademlia.GetLocalIP()
	fmt.Println("IP", localIP)

	kademliaNode := kademlia.NewNetwork(kademlia.NewKademlia())

	kademlia.Listen(localIP, 8080, kademliaNode)
}
