package pcaphandler_test

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/stretchr/testify/assert"
)

func TestConverFilePcap(t *testing.T) {
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
}

func TestReadPcapFile(t *testing.T) {
	handle, err := pcap.OpenOffline("../test_files/1616152209_2021_03_19____14_10_09_51841.tdp")
	assert.NoError(t, err)

	packets := gopacket.NewPacketSource(handle, handle.LinkType()).Packets()
	for pkt := range packets {
		fmt.Println(string(pkt.Data()))
	}

	handle.Close()

	assert.True(t, true)
}
