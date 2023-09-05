package main

import (
	"fmt"
	"time"

	"github.com/Domilz/D7024E/pkg/kademlia"
)

func main() {
	fmt.Println("Starting...")

	go kademlia.Listen("", 8080)

	kademliaNode := kademlia.NewNetwork(kademlia.NewKademlia())

	// bootstrapID := kademlia.NewKademliaID(DefaultBootstrapInput)

	// bootstrapNodeHostname := os.Getenv("BOOTSTRAP_NODE_IP")
	// ips, err := net.LookupIP(bootstrapNodeHostname)
	// if err != nil {
	// 	fmt.Println("Error", err)
	// }
	// ip := ips[0].String() + ":" + os.Getenv("NODE_PORT")

	// bootstrapContact := kademlia.NewContact(, ip)
	// kademliaNode.Kademlia.RoutingTable.AddContact(bootstrapContact)

	for {
		time.Sleep(2 * time.Second)
		kademliaNode.SendPingMessage()
	}

}
