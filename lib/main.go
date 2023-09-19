package main

import (
	"fmt"

	"github.com/Domilz/D7024E/pkg/cli"
	"github.com/Domilz/D7024E/pkg/kademlia"
)

func main() {
	fmt.Println("Starting...")

	localIP := kademlia.GetLocalIP()
	fmt.Println("IP", localIP)

	kademliaNode := kademlia.NewKademlia()

	go kademlia.Listen(localIP, "8080", kademliaNode)

	cli.CLI(kademliaNode)
}
