//go:build linux
// +build linux

package utils

import (
	"net"

	"github.com/artemis19/viz/pb"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	godpi "github.com/mushorg/go-dpi"
	"github.com/mushorg/go-dpi/types"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func SerializePacket(pkt gopacket.Packet) *pb.Packet {
	packet := &pb.Packet{}

	if ip4 := pkt.Layer(layers.LayerTypeIPv4); ip4 != nil {
		flow, _ := godpi.GetPacketFlow(pkt)
		result := godpi.ClassifyFlow(flow)
		// Setting default values to try and fix out of order packet fields
		packet.SrcAddr = "Source-IPAddress"
		packet.SrcPort = int32(99999)
		packet.DstAddr = "Dest-IPAddress"
		packet.DstPort = int32(99999)
		packet.Protocol = "N-A"

		packet.SrcAddr = net.IP(ip4.(*layers.IPv4).SrcIP).String()
		packet.DstAddr = net.IP(ip4.(*layers.IPv4).DstIP).String()

		if icmpLayer := pkt.Layer(layers.LayerTypeICMPv4); icmpLayer != nil {
			packet.Protocol = "ICMP"
		}

		if tcpLayer := pkt.Layer(layers.LayerTypeTCP); tcpLayer != nil {
			packet.Protocol = "TCP"
			packet.SrcPort = int32(tcpLayer.(*layers.TCP).SrcPort)
			packet.DstPort = int32(tcpLayer.(*layers.TCP).DstPort)
		}

		if udpLayer := pkt.Layer(layers.LayerTypeUDP); udpLayer != nil {
			packet.Protocol = "UDP"
			packet.SrcPort = int32(udpLayer.(*layers.UDP).SrcPort)
			packet.DstPort = int32(udpLayer.(*layers.UDP).DstPort)
		}
		if result.Protocol != types.Unknown {
			packet.Protocol = string(result.Protocol)
		}

		// Set default packet size to 0 if not specified
		packet.PktSize = int32(0)
		packet.PktSize = int32(pkt.Metadata().Length)
		// Set default timestamp
		packet.Timestamp = timestamppb.Now()
		packet.Timestamp = timestamppb.New(pkt.Metadata().Timestamp)
		// Set default value for machineID
		packet.MachineID = "MACHINE-ID"
		if MachineID != "" {
			packet.MachineID = MachineID
		}
	} else {
		// IPv6 packet ... ignoring it!
		return nil
	}

	return packet
}
