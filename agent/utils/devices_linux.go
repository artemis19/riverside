//go:build linux
// +build linux

package utils

import (
	"fmt"
	"net"

	"github.com/artemis19/viz/pb"
	"github.com/pkg/errors"
)

func GetNetInterfaceNames() ([]string, error) {
	// Make empty slice to store interface names
	var interfaceNames []string

	// Gets interfaces for linux OS
	interfaces, err := net.Interfaces()
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
	// Make list of interfaces
	interfaces := make([]*pb.NetworkInterface, 0)

	// Gets interfaces for Linux OS
	interfaceList, err := net.Interfaces()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to retrieve interfaces %v\n", err))
	}

	// Loop through them all
	for _, device := range interfaceList {
		netInt := &pb.NetworkInterface{}
		netInt.Name = device.Name
		addrs, err := device.Addrs()
		if err != nil {
			return nil, errors.New("Failed to retrieve network interfaces.")
		}
		for _, addr := range addrs {
			// Verify that it's an IPv4 address
			if ipv4Addr := addr.(*net.IPNet).IP.To4(); ipv4Addr != nil {
				netInt.IPAddress = addr.(*net.IPNet).IP.To4().String()
			}
		}
		// Add interface if we see IP address
		if netInt.IPAddress != "" {
			interfaces = append(interfaces, netInt)
		}
	}

	return interfaces, nil
}
