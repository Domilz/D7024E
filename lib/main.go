package main

import (
	"fmt"
	"time"

	"github.com/Domilz/D7024E/pkg/kademlia"
)

func main() {
	fmt.Println("Starting...")

	go kademlia.Listen("", 8080)

	net := kademlia.NewNetwork(kademlia.NewKademlia())
	for {
		time.Sleep(2 * time.Second)
		net.SendPingMessage()
	}

}
