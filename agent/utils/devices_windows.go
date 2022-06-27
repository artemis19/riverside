//go:build windows
// +build windows

package utils

import (
	"fmt"
	"github.com/artemis19/viz/pb"
	"github.com/google/gopacket/pcap"
	"github.com/pkg/errors"
)

func GetNetInterfaceNames() ([]string, error) {
	// Make empty slice to store interface names
	var interfaceNames []string

	// Gets interfaces for linux OS
	interfaces, err := pcap.FindAllDevs()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to retrieve interfaces %v\n", err))
	}

	// Loop through them all
	for _, device := range interfaces {
		interfaceNames = append(interfaceNames, device.Name)
	}

	return interfaceNames, nil
}

func GetInterfaceProtobufs() ([]*pb.NetworkInterface, error) {
	// Make map of protobuf interfaces
	interfaces := make([]*pb.NetworkInterface, 0)

	// Gets interfaces for Windows OS
	devices, err := pcap.FindAllDevs()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to retrieve interfaces %v\n", err))
	}

	for _, device := range devices {
		netInt := &pb.NetworkInterface{}
		// For Windows, use description for name
		netInt.Name = device.Description
		for _, address := range device.Addresses {
			// Verify that it's an IPv4 address
			if ipv4Addr := address.IP.To4(); ipv4Addr != nil {
				netInt.IPAddress = address.IP.To4().String()
			}
		}
		// Add interface if we see IP address
		if netInt.IPAddress != "" {
			interfaces = append(interfaces, netInt)
		}
	}
	return interfaces, nil
}
