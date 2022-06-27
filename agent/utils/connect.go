package utils

import (
	"errors"
	"fmt"
	"github.com/artemis19/viz/pb"
	"github.com/denisbrodbeck/machineid"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var MachineID string

func CheckEmptyServer() error {
	if viper.Get("serverAddress") == "" {
		e := ("Empty server address! Please supply in config file or arguments.\n")
		return errors.New(e)
	}
	return nil
}

func ServerConnect(address string) (pb.VizClient, error) {
	// Set up a connection to the server
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		e := fmt.Sprintf("Could not connect: %v", err)
		return nil, errors.New(e)
	}
	client := pb.NewVizClient(conn)

	return client, nil
}

// Generate unique MachineID for host
func GetMachineID() (string, error) {
	if MachineID != "" {
		return MachineID, nil
	}
	var err error
	// Creates unique, hashed machineID
	MachineID, err = machineid.ProtectedID("myVizID")
	if err != nil {
		return "", errors.New(fmt.Sprintf("Failed get machine ID %v\n", err))
	}
	return MachineID, nil
}

func OpenPacketStream(client pb.VizClient) (pb.Viz_CollectClient, error, context.CancelFunc) {
	// Open stream connection from client to server
	ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	stream, err := client.Collect(ctx)

	if err != nil {
		e := fmt.Sprintf("Could not connect to collection server... %v", err)
		cancel()
		return nil, errors.New(e), cancel
	}

	return stream, nil, cancel
}

func SendPacket(stream pb.Viz_CollectClient, packet *pb.Packet) error {
	// Send packets from client to server
	err := stream.Send(packet)
	if err != nil {
		e := fmt.Sprintf("Error sending packets... %v", err)
		return errors.New(e)
	}

	return nil
}

func ClosePacketStream(stream pb.Viz_CollectClient, cancel context.CancelFunc) (*pb.Reply, error) {
	// Close the connection and send server reply
	reply, err := stream.CloseAndRecv()
	fmt.Println("Closing packet stream...")
	if err != nil {
		cancel()
		return nil, err
	}
	cancel()
	return reply, nil
}
