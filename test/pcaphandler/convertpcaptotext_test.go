package pcaphandler_test

import (
	"os"
	"testing"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/stretchr/testify/assert"
)

func TestConvertPcapToText(t *testing.T) {
	handle, err := pcap.OpenOffline("../test_files/test_pcap_file.pcap")
	assert.NoError(t, err)

	f, err := os.OpenFile("../test_files/test_pcap_file.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	assert.NoError(t, err)

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	packets := packetSource.Packets()

	for packet := range packets {
		_, err := f.WriteString(packet.String())
		assert.NoError(t, err)

		appLayer := packet.ApplicationLayer()

	}

	f.Close()
	handle.Close()
}
