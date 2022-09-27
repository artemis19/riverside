package cmd

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/artemis19/viz/agent/config"
	"github.com/artemis19/viz/agent/host"
	"github.com/artemis19/viz/agent/utils"
	"github.com/artemis19/viz/pb"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
)

func ListenOnInterface(interfaceName string, snapLen int32, serverStream pb.Viz_CollectClient, extraFilter string) {
	log.Printf("Tasked to listen on interface %v\n", interfaceName)

	handle, err := pcap.OpenLive(interfaceName, int32(snapLen), true, pcap.BlockForever)
	if err != nil {
		log.Fatalf("Failed to listen on interface %v with error %v\n", interfaceName, err)
	}
	defer handle.Close()

	// Set automatic filter to not duplicate server traffic in packets
	serverPieces := strings.Split(viper.Get("serverAddress").(string), ":")
	serverHost, serverPort := serverPieces[0], serverPieces[1]
	filter := fmt.Sprintf("not (host %v and port %v)", serverHost, serverPort)
	if !(viper.Get("filter").(string) == "") {
		filter += fmt.Sprintf(" and (%v)", viper.Get("filter").(string))
	}

	// Add in a suppled extra filter to ignore the server's database (passed from agent check-in)
	filter += fmt.Sprintf(" and (%v)", extraFilter)

	// Set filters for packet capture
	if err := handle.SetBPFFilter(filter); err != nil {
		log.Fatalf("Failed to set filter %v with error %v\n", viper.Get("filter"), err)
	}

	// Create packet data source
	packets := gopacket.NewPacketSource(handle, handle.LinkType()).Packets()

	fmt.Printf("Now listening on interface %v with filter %v\n", interfaceName, filter)

	// Loop through packets as we get them
	for pkt := range packets {
		pbPacket := utils.SerializePacket(pkt)

		// Added since we do not support IPv6 packets which will return nil
		if pbPacket != nil {
			log.Printf("%v -> %v", interfaceName, pbPacket)
			err = utils.SendPacket(serverStream, pbPacket)
			if err != nil {
				log.Fatalf("Error sending packets to server, %v, with error %v\n", serverHost, err)
			}
		}
	}
}

func Listen(cmd *cobra.Command, args []string) {
	// Set defaults & start logging
	config.InitializeConfig()

	log.Printf("Tasked to listen...\n")
	defer log.Printf("Done listening.\n")

	var serverAddress string = viper.Get("serverAddress").(string)

	err := utils.CheckEmptyServer()
	if err != nil {
		log.Fatalf("Failed to ping server %v: %v\n", serverAddress, err)
		return
	}

	// Connect to the server
	client, err := utils.ServerConnect(serverAddress)
	if err != nil {
		log.Fatalf("Could not connect to server %v\n", err)
	}

	// Get agent host data
	host, err := host.NewHost()
	if err != nil {
		log.Fatalf("Failed to gather host data for this agent: %v\n", err)
		return
	}

	// Send check-in request
	checkInResponse, err := client.CheckIn(context.Background(), host)
	if err != nil {
		log.Fatalf("Failed to send check-in request to the server: %v\n", err)
		return
	}

	log.Printf("Sucessfully checked in with response %v!\n", checkInResponse)
	extraFilter := checkInResponse.Message

	serverStream, err, cancel := utils.OpenPacketStream(client)
	if err != nil {
		log.Fatalf("Could not connect to collection server... %v\n", err)
	}

	defer utils.ClosePacketStream(serverStream, cancel)

	// Properly handle the packet length
	snapLen, err := strconv.Atoi(viper.Get("snapLength").(string))
	if err != nil {
		log.Fatalf("Erorr converting string to integer.\n")
	}

	// Create wait group since we need to wait for each interface's goroutine to finish
	var wg sync.WaitGroup

	interfaces, err := utils.GetNetInterfaceNames()
	if err != nil {
		log.Fatalf("Failed to find interfaces %v\n.", err)
	}

	for _, interfaceName := range interfaces {
		if strings.HasPrefix(interfaceName, "veth") || strings.HasPrefix(interfaceName, "docker") || strings.HasPrefix(interfaceName, "lo") {
			// Ignore virtual, docker, and loopback interfaces
			continue
		}
		// Add every interface to our group
		wg.Add(1)
		defer wg.Done()

		// Goroutine so all of the device's interfaces can listen
		go ListenOnInterface(interfaceName, int32(snapLen), serverStream, extraFilter)

	}

	// Let every interface listen in the group
	wg.Wait()
}
