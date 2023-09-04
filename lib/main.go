package main

import (
	"github.com/Domilz/D7024E/pkg/kademlia"
)

func main() {
	// x := kademlia.NewKademlia()
	kademlia.Listen("", 8080)
}
