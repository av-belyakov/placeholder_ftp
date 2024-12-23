package pcaphandler_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/stretchr/testify/assert"
)

/*func TestConverFilePcap(t *testing.T) {
	//handler, err := pcap.OpenOffline("../test_files/1616152209_2021_03_19____14_10_09_51841.tdp")
	f, err := os.OpenFile("../test_files/1616152209_2021_03_19____14_10_09_51841.tdp", os.O_RDONLY, os.ModePerm)
	assert.NoError(t, err)

	var wr bytes.Buffer
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		_, err := wr.Write(sc.Bytes())
		assert.NoError(t, err)
	}

	gp := gopacket.NewPacket(wr.Bytes(), layers.LayerTypeEthernet, gopacket.Default)

	fmt.Println("RESULT", gp.String())

	f.Close()

	assert.True(t, true)
}*/

func getPacketTypeIPv4Info(packet gopacket.Packet) {
	ipLayer := packet.Layer(layers.LayerTypeIPv4)

	for ipLayer != nil {
		fmt.Println("IPv4 layer detected.")
		ip, _ := ipLayer.(*layers.IPv4)
		fmt.Printf("From host %d.%d.%d.%d to %d.%d.%d.%d\n", ip.SrcIP[0], ip.SrcIP[1], ip.SrcIP[2], ip.SrcIP[3], ip.DstIP[0], ip.DstIP[1], ip.DstIP[2], ip.DstIP[3])
	}
}

func getPacketTransportInfo(packet gopacket.Packet) {
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

func getPacketApplicationLayer(packet gopacket.Packet) {
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

func TestReadPcapFile(t *testing.T) {
	handle, err := pcap.OpenOffline("../test_files/1616152209_2021_03_19____14_10_09_51841.tdp")
	assert.NoError(t, err)

	//packets := gopacket.NewPacketSource(handle, handle.LinkType()).Packets()
	//for pkt := range packets {
	//	fmt.Println(string(pkt.Data()))
	//}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		fmt.Println("=== Packet Info ===")

		getPacketTransportInfo(packet)
		getPacketTypeIPv4Info(packet)
		getPacketApplicationLayer(packet)

		fmt.Println("===================")
	}

	handle.Close()

	assert.True(t, true)
}
