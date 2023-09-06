package main

import (
	"fmt"

	"github.com/Domilz/D7024E/pkg/kademlia"
)

func main() {
	fmt.Println("Starting...")

	kademlia.Listen("", 8080)
}
