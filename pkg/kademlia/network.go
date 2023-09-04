package kademlia

import (
	"net"
	"strconv"
)

type Network struct {
	kademlia Kademlia
}

func Listen(ip string, port int) {
	address := ip + ":" + strconv.Itoa(port)
	net.Dial("udp", address)

	// TODO
}

func (network *Network) SendPingMessage(contact *Contact) {
	// TODO

}

func (network *Network) SendFindContactMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}
