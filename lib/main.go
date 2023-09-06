package main

import (
	"fmt"
	"net"
	"os"

	"github.com/Domilz/D7024E/pkg/kademlia"
)

func main() {
	fmt.Println("Starting...")

	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println(fmt.Errorf("error hostname lookup: %w", err))
	}

	localIP, err := net.LookupIP(hostname)
	if err != nil {
		fmt.Println(fmt.Errorf("error IP lookup: %w", err))
	}
	fmt.Println("IP", localIP[0])

	kademlia.Listen(string(localIP[0]), 8080)
}
