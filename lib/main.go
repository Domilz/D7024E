package main

import (
	"context"
	"fmt"
	"log"
	"os/exec"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func main() {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatal(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	// Loop through the containers and get their IP addresses
	for _, container := range containers {
		inspect, err := cli.ContainerInspect(context.Background(), container.ID)
		if err != nil {
			log.Fatal(err)
		}

		// Get the container's IP address from its network settings
		ip := inspect.NetworkSettings.Networks["bridge"].IPAddress
		fmt.Printf("Container ID: %s, IP Address: %s\n", container.ID[:10], ip)

		// Ping the container
		pingCommand := exec.Command("ping", "-c", "4", ip)
		pingOutput, err := pingCommand.CombinedOutput()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(pingOutput))
	}
}
