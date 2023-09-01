package d7024e

import (
	"strconv"
)

type Network struct {
}

func Listen(ip string, port int) {
	_ = ip + ":" + strconv.Itoa(port)
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
