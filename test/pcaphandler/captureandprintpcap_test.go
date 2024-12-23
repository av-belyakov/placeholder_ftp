package pcaphandler_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/stretchr/testify/assert"
)

var (
	//device      string = "\\Device\\NPF_{C410B1B0-56DE-4CD5-BC7A-5A5ACAB7619F}"
	device      string = "enp2s0"
	snaplen     int32  = 65536
	promiscuous bool   = true
	err         error
	timeout     time.Duration = -1 * time.Second
	handle      *pcap.Handle
)

func printEthInfo(packet gopacket.Packet) {
	ethernetLayer := packet.Layer(layers.LayerTypeEthernet)

	if ethernetLayer != nil {
		fmt.Println("Ethernet layer detected.")
		ethernetPacket, _ := ethernetLayer.(*layers.Ethernet)
		fmt.Println("Source MAC: ", ethernetPacket.SrcMAC)
		fmt.Println("Destination MAC: ", ethernetPacket.DstMAC)
		// Ethernet type is typically IPv4 but could be ARP or other
		fmt.Println("Ethernet type: ", ethernetPacket.EthernetType)
		fmt.Println()
	}
}

func printNetworkInfo(packet gopacket.Packet) {
	// Let's see if the packet is IP (even though the ether type told us)
	ipLayer := packet.Layer(layers.LayerTypeIPv4)

	if ipLayer != nil {
		fmt.Println("IPv4 layer detected.")
		ip, _ := ipLayer.(*layers.IPv4)
		// IP layer variables:
		// Version (Either 4 or 6)
		// IHL (IP Header Length in 32-bit words)
		// TOS, Length, Id, Flags, FragOffset, TTL, Protocol (TCP?),
		// Checksum, SrcIP, DstIP
		fmt.Printf("From %s to %s\n", ip.SrcIP, ip.DstIP)
		fmt.Println("Protocol: ", ip.Protocol)
		fmt.Println()
	}
}
func printTransportInfo(packet gopacket.Packet) {
	tcpLayer := packet.Layer(layers.LayerTypeTCP)

	if tcpLayer != nil {
		fmt.Println("TCP layer detected.")
		tcp, _ := tcpLayer.(*layers.TCP)
		// TCP layer variables:
		// SrcPort, DstPort, Seq, Ack, DataOffset, Window, Checksum, Urgent
		// Bool flags: FIN, SYN, RST, PSH, ACK, URG, ECE, CWR, NS
		fmt.Printf("From port %d to %d\n", tcp.SrcPort, tcp.DstPort)
		fmt.Println("Sequence number: ", tcp.Seq)
		fmt.Println()
	}
}
func printApplicationInfo(packet gopacket.Packet) {
	applicationLayer := packet.ApplicationLayer()

	if applicationLayer != nil {
		fmt.Println("Application layer/Payload found.")
		fmt.Printf("%s\n", applicationLayer.LayerType())

		// Search for a string inside the payload
		if strings.Contains(string(applicationLayer.Payload()), "HTTP") {
			fmt.Println("HTTP found!")
		}
	}
}

func printLayer(packet gopacket.Packet) {
	fmt.Println("All packet layers:")

	for _, layer := range packet.Layers() {
		fmt.Println("- ", layer.LayerType())
	}
}

func printErrInfo(packet gopacket.Packet) {
	// Check for errors
	if err := packet.ErrorLayer(); err != nil {
		fmt.Println("Error decoding some part of the packet:", err)
	}
}

func TestCaptureNetworkTraffik(t *testing.T) {
	devices, err := pcap.FindAllDevs()
	assert.NoError(t, err)

	/*fmt.Println("Devices found:")
	for _, device := range devices {
		fmt.Println("\nName: ", device.Name)
		fmt.Println("Description: ", device.Description)
		fmt.Println("Devices flag: ", device.Flags)
		for _, address := range device.Addresses {
			fmt.Println("- IP address: ", address.IP)
			fmt.Println("- Subnet mask: ", address.Netmask)
			fmt.Println("- Broadaddr:  ", address.Broadaddr)
			fmt.Println("- P2P:  ", address.P2P)
		}
	}*/

	assert.Greater(t, len(devices), 0)

	// Open device
	handle, err = pcap.OpenLive(device, snaplen, promiscuous, timeout)
	assert.NoError(t, err)

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		fmt.Println("=========================")

		printEthInfo(packet)
		printNetworkInfo(packet)
		printTransportInfo(packet)
		printApplicationInfo(packet)
		printLayer(packet)
		printErrInfo(packet)
	}

	handle.Close()
}
