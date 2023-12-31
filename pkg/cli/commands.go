package cli

import (
	"encoding/hex"
	"fmt"

	"github.com/Domilz/D7024E/pkg/kademlia"
)

func get(args []string, kademlia *kademlia.Kademlia) {
	if len(args) == 1 {
		decoded, _ := hex.DecodeString(args[0])
		if len(decoded) < 20 {
			fmt.Println("Insufficient argument")
			return
		}

		data, err := kademlia.LookupData(args[0])
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("Contact %v had the content: %v", data.Contact.ID, data.Value)
		}
	} else {
		fmt.Println("Bad argument")
	}
}

func put(args []string, kademlia *kademlia.Kademlia) {
	if len(args) == 1 {
		hash, err := kademlia.Store([]byte(args[0]))
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Hash :=", hash.String())
		}
	} else {
		fmt.Println("Bad argument")
	}
}

func help() string {
	return `
Kademlia CLI to execute commands

USAGE:
	<command> [arguments]
	
The commands are:
	exit      Terminates and exit
	put       Appends content to network
	get       Retrieves content
	help      Show help
`
}
