package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Domilz/D7024E/pkg/kademlia"
)

var (
	in = os.Stdin
)

func processInput(input string, kademlia *kademlia.Kademlia) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return
	}

	command := parts[0]
	args := parts[1:]

	switch command {
	case "put":
		put(args, kademlia)
	case "get":
		get(args, kademlia)
	case "exit":
		fmt.Println("Exiting node")
		os.Exit(0)
	case "help":
		fmt.Println(help())
	default:
		fmt.Println("Unknown command:", command)
		fmt.Println(help())
	}
}

func CLI(kademlia *kademlia.Kademlia) {
	fmt.Println("Starting CLI...")
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		input := scanner.Text()

		processInput(input, kademlia)
	}

	if scanner.Err() != nil {
		fmt.Fprintln(os.Stderr, "Error reading input:", scanner.Err())
	}

}
