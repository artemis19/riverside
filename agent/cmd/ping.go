package cmd

import (
	// "fmt"
	"github.com/artemis19/viz/agent/config"
	"github.com/artemis19/viz/agent/utils"
	"github.com/artemis19/viz/pb"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

func Ping(cmd *cobra.Command, args []string) {
	// Set defaults & start logging
	config.InitializeConfig()

	var serverAddress string = viper.Get("serverAddress").(string)

	log.Printf("Tasked to ping server address %v...\n", serverAddress)
	defer log.Printf("Done pinging server address.\n")

	err := utils.CheckEmptyServer()
	if err != nil {
		log.Fatalf("Failed to ping server %v: %v\n", serverAddress, err)
		return
	}

	client, err := utils.ServerConnect(serverAddress)
	if err != nil {
		log.Fatalf("Could not connect to server %v\n", err)
	}

	stream, err, cancel := utils.OpenPacketStream(client)
	if err != nil {
		log.Fatalf("Could not connect to collection server... %v", err)
	}

	err = utils.SendPacket(stream, &pb.Packet{})
	if err != nil {
		log.Fatalf("Error sending packets... %v", err)
	}

	reply, err := utils.ClosePacketStream(stream, cancel)
	if err != nil {
		log.Fatalf("Error during stream... %v", err)
	}

	log.Printf("Received reply from server %v\n", reply)
}
