package main

import (
	"fmt"
	"os"

	"github.com/Domilz/D7024E/pkg/cli"
	"github.com/Domilz/D7024E/pkg/kademlia"
)

func main() {
	fmt.Println("Starting...")

	localIP := kademlia.GetLocalIP()
	fmt.Println("IP", localIP)

	kademliaNode := kademlia.NewKademlia()

	isBootstrap := os.Getenv("ISBOOTSTRAP")

	if isBootstrap == "0" {
		kademliaNode.JoinNetwork()
	}

	go kademlia.Listen(localIP, "8080", kademliaNode)

	cli.CLI(kademliaNode)
}
