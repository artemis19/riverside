package host

import (
	"errors"
	"fmt"
	"github.com/artemis19/viz/agent/utils"
	"github.com/artemis19/viz/pb"
	"os"
	"runtime"
)

func NewHost() (*pb.Host, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to retrieve hostname: %v\n", err))
	}

	// Get hashed machine ID
	MachineID, err := utils.GetMachineID()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed get machine ID %v\n", err))
	}

	// Built out for Windows or Linux in utils/devices
	interfaces, err := utils.GetInterfaceProtobufs()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to retrieve network interfaces %v\n", err))
	}

	host := &pb.Host{
		Hostname:     hostname,
		OS:           runtime.GOOS,
		Architecture: runtime.GOARCH,
		MachineID:    MachineID,
		Interfaces:   interfaces,
	}

	return host, nil
}
