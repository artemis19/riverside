package rpc

import (
	"fmt"

	"github.com/artemis19/viz/pb"
	"github.com/artemis19/viz/rpc/sqltime"
	"github.com/artemis19/viz/server/database"
	"github.com/artemis19/viz/server/websocket"

	"crypto/sha256"
	"encoding/hex"

	"io"
	"log"
	"net"
	"time"

	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type VizServer struct {
	pb.UnimplementedVizServer
}

// Sets up initial server connection based on arguments supplied
func RunGRPCServer() {
	listenPort := fmt.Sprintf(":%v", viper.Get("port"))
	listener, err := net.Listen("tcp", listenPort)

	if err != nil {
		log.Fatalf("Server failed to listen %v\n", err)
	}

	s := grpc.NewServer()
	pb.RegisterVizServer(s, &VizServer{})
	log.Printf("Server listening at %v\n", listener.Addr())

	if err := s.Serve(listener); err != nil {
		log.Fatalf("Failed to serve %v\n", err)
	}
}

func (s *VizServer) Collect(stream pb.Viz_CollectServer) error {
	var packetCount int32
	startTime := time.Now()

	for {
		packet, err := stream.Recv()
		log.Printf("Received packet %v", packet)

		if packet != nil {
			remoteHost := HandleRemoteHost(packet)
			log.Printf("Communicating with identified remote host '%v'", remoteHost)
			HandleNetFlow(packet, remoteHost)
		}

		if err == io.EOF {
			endTime := time.Now()
			message := fmt.Sprintf("start: %v - end: %v - received %v packets", startTime, endTime, packetCount)
			log.Printf("Message: %v\n", message)
			return stream.SendAndClose(&pb.Reply{
				Message: message,
			})
		}

		if err != nil {
			log.Printf("Error while collecting on server: %v.\n", err)
			return err
		}
		packetCount++
	}
}

/*
Agent - sends CheckIn
Server - When receiving CheckIn, create/store new host
Client - Update front-end to display new host on network map
*/

func (s *VizServer) CheckIn(ctx context.Context, host *pb.Host) (*pb.Reply, error) {
	log.Printf("%v checked in!\n", host)

	dbHost := &database.Host{
		Hostname:     host.Hostname,
		OS:           host.OS,
		Architecture: host.Architecture,
		MachineID:    host.MachineID,
	}

	var alreadyExists database.Host
	// Check if host already exists in database
	query := database.DB.Where("machine_id = ?", host.MachineID).First(&alreadyExists)
	if query.Error != nil {
		log.Printf("This machine is not yet found in the database, adding entry now...")
		for _, device := range host.Interfaces {
			dbHost.Interfaces = append(dbHost.Interfaces, database.NetworkInterface{
				Name:      device.Name,
				IPAddress: device.IPAddress,
			})
		}
		// Create host object in DB w/network interfaces
		database.DB.Create(&dbHost)
		sendWSAgent(*dbHost)
	} else {
		log.Printf("This machine is already present in the database!")
		// Update host fields (if they changed)
		database.DB.Model(&alreadyExists).Updates(database.Host{
			Hostname:     host.Hostname,
			OS:           host.OS,
			Architecture: host.Architecture,
		})

		// First, make sure network interfaces present in the database match the host
		// Remove if necessary
		var dbInterfaces []database.NetworkInterface
		database.DB.Where("host_id = ?", host.MachineID).Find(&dbInterfaces)
		for _, eachInterface := range dbInterfaces {
			var found bool
			found = false
			for _, device := range host.Interfaces {
				if eachInterface.Name == device.Name {
					found = true
					break
				}
			}
			if found {
				continue
			} else {
				database.DB.Unscoped().Delete(&eachInterface)
			}
		}

		// Loop through network interfaces to see if they've changed
		// Update ones that exist, add new ones if necessary
		for _, device := range host.Interfaces {
			var db_interface database.NetworkInterface
			query := database.DB.Where("host_id = ? AND name = ?", host.MachineID, device.Name).First(&db_interface)
			if query.Error != nil {
				database.DB.Create(&database.NetworkInterface{
					HostID:    host.MachineID,
					Name:      device.Name,
					IPAddress: device.IPAddress,
				})
			} else {
				database.DB.Model(&db_interface).Updates(database.NetworkInterface{
					Name:      device.Name,
					IPAddress: device.IPAddress,
				})
			}
		}
	}

	// Tell agent what the database filter is so we don't duplicate traffic
	dbAddress := viper.Get("dbAddress").(string)
	dbPort := viper.Get("dbPort").(string)
	extraFilter := fmt.Sprintf("not (host %v and port %v)", dbAddress, dbPort)
	return &pb.Reply{Message: extraFilter}, nil
}

// From context of a packet, determine whether this is a new remote host that needs to be added to the database
func HandleRemoteHost(packet *pb.Packet) string {
	// If SrcAdr or DstAdr is not in RemoteHosts table or if it is not a known machineId in network interfaces table,
	// then this is a new remote host that needs to be added

	log.Printf("Checking packet for remote hosts...")
	for _, eachIP := range []string{packet.SrcAddr, packet.DstAddr} {
		if eachIP == "" {
			continue
		}
		var seenRemoteHostBefore database.RemoteHost
		resultRemote := database.DB.Where("ip_address = ?", eachIP).Find(&seenRemoteHostBefore)
		if resultRemote.RowsAffected == 0 {
			log.Printf("IP address '%v' is not in RemoteHosts table!", eachIP)

			// Get network interfaces by machine ID for
			var networkInterfaces []database.NetworkInterface
			resultHost := database.DB.Where("ip_address = ?", eachIP).Find(&networkInterfaces)
			//  log.Printf("Found %v for '%v'", networkInterfaces, eachIP)

			// Check if src IP is an agent host
			if resultHost.RowsAffected == 0 {
				log.Printf("Packet IP address '%v' is not in Hosts table either... adding to RemoteHosts table...", eachIP)
				// Hash remote host IP to match machineID
				hasher := sha256.New()
				hasher.Write([]byte(eachIP))
				eachIPHash := hex.EncodeToString(hasher.Sum(nil))

				database.DB.Create(&database.RemoteHost{
					RemoteHostID: eachIPHash,
					IPAddress:    eachIP,
				})
				return eachIP
			} else {
				log.Printf("Remote host IP '%v' is in the NetworkInterfaces table, indicaing internal traffic.", eachIP)
			}
		} else {
			log.Printf("Remote host IP '%v' is in the RemoteHosts table...", eachIP)
			return eachIP
		}
	}
	// No remote host was identified, internal traffic
	return ""
}

// Convenience function to check if a remote host is in our map of previous remote host IDs
func isInSlice(remoteHost string, prevRemoteIDs []string) bool {
	for _, element := range prevRemoteIDs {
		if element == remoteHost {
			return true
		}
	}
	return false
}

// Tell websockets new netflow exists
func sendWSNetFlow(data database.NetFlow) {
	// dispatch, err := json.Marshal(flow)
	// if err != nil {
	// 	log.Printf("Failed to marshal flow data: %v", err)
	// }
	// wsDispatch, err := websocket.SendAll(dispatch, "netflow")
	// if err != nil {
	// 	log.Printf("Failed to prep websocket data for SendAll(): %v", err)
	// }
	websocket.WSNetflowWriter <- data
}

// Tell websockets new host exists
func sendWSAgent(data database.Host) {
	// dispatch, err := json.Marshal(flow)
	// if err != nil {
	// 	log.Printf("Failed to marshal flow data: %v", err)
	// }
	// wsDispatch, err := websocket.SendAll(dispatch, "host")
	// if err != nil {
	// 	log.Printf("Failed to prep websocket data for SendAll(): %v", err)
	// }
	websocket.WSHostWriter <- data
}

func HandleNetFlow(packet *pb.Packet, remoteHostIP string) {
	log.Printf("Netflow sees packet: %v", packet)
	// From the perspective of the agent who initiates the traffic
	var remoteHostID, srcIP, dstIP, direction string
	var srcPort, dstPort int
	// For calculating traffic in our "window" to incorporate the timeline and calculate throughput
	// Set current packet time and the previous second
	packetTime := packet.Timestamp.AsTime()
	lastSecond := packetTime.Add(time.Second * -2)
	// lastSecond = lastSecond.Add(time.Millisecond * -500)

	// Keep track of previous entries within past second in netflow table
	var prevEntries []database.NetFlow
	prevEntriesResult := database.DB.Where("host_id = ? AND end_time BETWEEN ? AND ?", packet.MachineID, sqltime.New(lastSecond), sqltime.New(packetTime)).Find(&prevEntries)

	// Make map of previous communications to look at previous hosts our initial host communicated with
	var prevRemoteIDs []string
	prevComms := make(map[string][]database.NetFlow, 0)

	// Check if traffic has occurred within the last second on the agent host
	if prevEntriesResult.RowsAffected != 0 {
		for _, entry := range prevEntries {
			remoteHost := entry.RemoteHostID
			if !isInSlice(remoteHost, prevRemoteIDs) {
				prevRemoteIDs = append(prevRemoteIDs, remoteHost)
				prevComms[remoteHost] = append(prevComms[remoteHost], entry)
			}
		}
	}

	log.Printf("Netflow thinks the remote host IP is '%v'", remoteHostIP)

	if remoteHostIP == "" {
		// This is for internal <-> internal traffic
		log.Printf("Remote host IP is empty... assuming internal <-> internal traffic.")
		// Looking for agent-installed host
		var networkInterfaces database.NetworkInterface
		result := database.DB.Where("ip_address = ? AND host_id = ?", packet.SrcAddr, packet.MachineID).Find(&networkInterfaces)
		log.Printf("Netflow: Found packet: %v", packet)
		if result.RowsAffected == 0 {
			// For purposes of "direction" in our database, we are always looking
			// at this from the perspective of the agent host who initiates the traffic
			// Make agent host (DstAddr) who is receiving the traffic the "remote" host, so direction is FROM
			var fromInternalNetworkInterfaces database.NetworkInterface
			fromResult := database.DB.Where("ip_address = ?", packet.SrcAddr).Find(&fromInternalNetworkInterfaces)
			log.Printf("Netflow FROM: Found %v for '%v'", fromInternalNetworkInterfaces, packet.SrcAddr)
			if fromResult.RowsAffected == 0 {
				// This shouldn't happen
				log.Printf("Netflow unable to find IP address '%v' from internal agents!", packet.DstAddr)
			} else {
				remoteHostID = fromInternalNetworkInterfaces.HostID
				dstIP = packet.SrcAddr
				srcIP = packet.DstAddr
				dstPort = int(packet.SrcPort)
				srcPort = int(packet.DstPort)
				direction = "from"
				log.Printf("Netflow: FROM perspective with remoteHostID '%v'", remoteHostID)

				// Bidirectional check
				if isInSlice(remoteHostID, prevRemoteIDs) {
					log.Printf("Netflow: This agent was talking to a previously identified host.")
					for _, previous := range prevComms[remoteHostID] {
						// Verify this is the same communication
						if previous.Protocol == packet.Protocol && (previous.SrcAddress == packet.DstAddr || previous.SrcAddress == packet.SrcAddr) && (int32(previous.SrcPort) == packet.DstPort || int32(previous.SrcPort) == packet.SrcPort) {
							if previous.Direction != direction {
								previous.Direction = "bidirectional"
							}
							// Increment throughput
							previous.Throughput += 1
							// Keep track of end time for continued bidirectional comms
							previous.EndTime = sqltime.New(packet.Timestamp.AsTime())
							// Save net flow to the database
							database.DB.Save(&previous)
							sendWSNetFlow(previous)
						}
					}
					return
				}
			}
		} else {
			var internalNetworkInterfaces database.NetworkInterface
			result := database.DB.Where("ip_address = ?", packet.DstAddr).Find(&internalNetworkInterfaces)
			log.Printf("Netflow TO: Found %v for '%v'", internalNetworkInterfaces, packet.DstAddr)
			if result.RowsAffected == 0 {
				// This shouldn't happen
				log.Printf("Netflow unable to find IP address '%v' from internal agents!", packet.DstAddr)
			} else {
				// These are now flipped since the direction is TO
				remoteHostID = internalNetworkInterfaces.HostID
				dstIP = packet.DstAddr
				srcIP = packet.SrcAddr
				dstPort = int(packet.DstPort)
				srcPort = int(packet.SrcPort)
				direction = "to"
				log.Printf("Netflow: TO perspective with remoteHostID '%v'", remoteHostID)

				// Temporary check
				if isInSlice(remoteHostID, prevRemoteIDs) {
					log.Printf("Netflow: This agent was talking to a previously identified host.")
					for _, previous := range prevComms[remoteHostID] {
						// Verify this is the same communication
						if previous.Protocol == packet.Protocol && (previous.SrcAddress == packet.DstAddr || previous.SrcAddress == packet.SrcAddr) && (int32(previous.SrcPort) == packet.DstPort || int32(previous.SrcPort) == packet.SrcPort) {
							if previous.Direction != direction {
								previous.Direction = "bidirectional"
							}
							// Increment throughput
							previous.Throughput += 1
							// Keep track of end time for continued bidirectional comms
							previous.EndTime = sqltime.New(packet.Timestamp.AsTime())
							// Save net flow to the database
							database.DB.Save(&previous)
							sendWSNetFlow(previous)
						}
					}
					return
				}
			}
		}
	} else {
		// This is for external <-> internal traffic
		log.Printf("Netflow has a remote host IP '%v'... assuming external <-> internal traffic.", remoteHostIP)
		hasher := sha256.New()
		hasher.Write([]byte(remoteHostIP))
		remoteHostID = hex.EncodeToString(hasher.Sum(nil))

		// Source address is agent's IP address since we assume the perspective of the internal agent host
		// Need IP then that is not the remoteHostIP
		srcIP = packet.SrcAddr
		dstIP = remoteHostIP
		if packet.SrcAddr == remoteHostIP {
			// Make the agent node the "source"
			srcIP = packet.DstAddr
		}

		// Same thing for source port
		srcPort = int(packet.SrcPort)
		dstPort = int(packet.DstPort)
		if packet.SrcAddr == remoteHostIP {
			// If opposite perspective, swap
			srcPort = int(packet.DstPort)
			dstPort = int(packet.SrcPort)
		}

		// Direction is again from agent's perspective
		// If A -> B, TO
		// If B -> A, FROM
		direction = "from"
		if packet.DstAddr == remoteHostIP {
			direction = "to"
		}

		if isInSlice(remoteHostID, prevRemoteIDs) {
			log.Printf("Netflow: This agent was talking to a previously identified host.")
			for _, previous := range prevComms[remoteHostID] {
				// Verify this is the same communication
				if previous.Protocol == packet.Protocol && (previous.SrcAddress == packet.DstAddr || previous.SrcAddress == packet.SrcAddr) && (int32(previous.SrcPort) == packet.DstPort || int32(previous.SrcPort) == packet.SrcPort) {
					if previous.Direction != direction {
						previous.Direction = "bidirectional"
					}
					// Increment throughput
					previous.Throughput += 1
					// Keep track of end time for continued bidirectional comms
					previous.EndTime = sqltime.New(packet.Timestamp.AsTime())
					// Save net flow to the database
					database.DB.Save(&previous)
					sendWSNetFlow(previous)
				}
			}
			return
		}
	}

	// Agent who initiates the traffic is the "source" here
	flow := &database.NetFlow{
		HostID:       packet.MachineID,
		RemoteHostID: remoteHostID,
		SrcAddress:   srcIP,
		DstAddress:   dstIP,
		SrcPort:      srcPort,
		DstPort:      dstPort,
		Throughput:   1,
		Direction:    direction,
		Protocol:     packet.Protocol,
		StartTime:    sqltime.New(packet.Timestamp.AsTime()),
		EndTime:      sqltime.New(packet.Timestamp.AsTime()),
	}
	database.DB.Create(&flow)
	sendWSNetFlow(*flow)
	log.Printf("Tracking netflow: %v", flow)
}
