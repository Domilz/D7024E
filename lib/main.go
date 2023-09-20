package main

import (
	"fmt"
	"os"
	"time"

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
		time.Sleep(5 * time.Second)
		kademliaNode.JoinNetwork()
	}

	port := os.Getenv("NODE_PORT")

	go kademlia.Listen(localIP, port, kademliaNode)

	cli.CLI(kademliaNode)
}
