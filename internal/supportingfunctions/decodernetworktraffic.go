package supportingfunctions

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
	"unicode"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcapgo"

	"github.com/av-belyakov/placeholder_ftp/cmd/commoninterfaces"
)

// NetworkTrafficDecoder декодировщик сетевого трафика
func NetworkTrafficDecoder(fileName string, fr, fw *os.File, logger commoninterfaces.Logger) error {
	r, err := pcapgo.NewReader(fr)
	if err != nil {
		return err
	}

	writer := bufio.NewWriter(fw)
	defer func() {
		if err == nil {
			err = writer.Flush()
		}
	}()

	_, err = writer.WriteString(fmt.Sprintf("Decoding file name: '%v', decoding time: %v\n", fileName, time.Now().Format(time.RFC3339)))
	if err != nil {
		return err
	}

	var (
		eth     layers.Ethernet
		ip4     layers.IPv4
		ip6     layers.IPv6
		tcp     layers.TCP
		udp     layers.UDP
		dns     layers.DNS
		ntp     layers.NTP
		tls     layers.TLS
		payload gopacket.Payload
	)

	decoded := make([]gopacket.LayerType, 10)
	parser := gopacket.NewDecodingLayerParser(layers.LayerTypeEthernet, &eth, &ip4, &ip6, &tcp, &udp, &dns, &ntp, &tls, &payload)

	boolToInt8 := func(v bool) int8 {
		if v {
			return 1
		}

		return 0
	}

	for {
		data, ci, e := r.ReadPacketData()
		if e != nil {
			if e == io.EOF {
				break
			}
		}

		e = parser.DecodeLayers(data, &decoded)
		if e != nil {
			continue
		}

		var errWrite error
		for _, layerType := range decoded {
			switch layerType {
			case layers.LayerTypeIPv6:
				_, errWrite = writer.WriteString(fmt.Sprintf("\n\n\n%v, packets length: %v\nIP6 %v -> %v\n", ci.Timestamp, ci.CaptureLength, ip6.SrcIP, ip6.DstIP))

			case layers.LayerTypeIPv4:
				_, errWrite = writer.WriteString(fmt.Sprintf("\n\n\n%v, packets length: %v\nIP4 %v -> %v\n", ci.Timestamp, ci.CaptureLength, ip4.SrcIP, ip4.DstIP))

			case layers.LayerTypeTCP:
				payloadSize := len(tcp.LayerPayload())
				if _, errWrite = writer.WriteString(fmt.Sprintf("TCP port %v -> %v\n", tcp.SrcPort, tcp.DstPort)); errWrite != nil {
					continue
				}

				fin := boolToInt8(tcp.FIN)
				syn := boolToInt8(tcp.SYN)
				rst := boolToInt8(tcp.RST)
				psh := boolToInt8(tcp.PSH)
				ack := boolToInt8(tcp.ACK)
				urg := boolToInt8(tcp.URG)

				if _, errWrite = writer.WriteString(fmt.Sprintf("Flags (FIN:'%v' SYN:'%v' RST:'%v' PSH:'%v' ACK:'%v' URG:'%v')\n", fin, syn, rst, psh, ack, urg)); errWrite != nil {
					continue
				}

				if payloadSize == 0 {
					continue
				}

				reader := bufio.NewReader(bytes.NewReader(tcp.LayerPayload()))
				reqHttp, errHttp := http.ReadRequest(reader)
				if errHttp == nil {
					proto := reqHttp.Proto
					method := reqHttp.Method
					//url := httpReq.URL //содержит целый тип, не только значение httpReq.RequestURI но и методы для парсинга запроса
					host := reqHttp.Host
					reqURI := reqHttp.RequestURI
					userAgent := reqHttp.Header.Get("User-Agent")
					contentType := reqHttp.Header.Get("Content-Type")

					reqHttp.Body.Close()

					if _, errWrite = writer.WriteString(fmt.Sprintf("\n%v %v %v\nContent-Type:%v\nHost:%v\nUser-Agent:%v\n", method, reqURI, proto, contentType, host, userAgent)); errWrite != nil {
						continue
					}
				} else {
					if strings.Contains(string(tcp.LayerPayload()), "HTTP/") {
						_, errWrite = writer.WriteString(fmt.Sprintf("\n%v\n", strings.TrimFunc(string(payload), func(r rune) bool {
							return unicode.IsSpace(r)
						})))
					}
				}

			case layers.LayerTypeUDP:
				_, errWrite = writer.WriteString(fmt.Sprintf("UDP port:%v -> %v\n", udp.SrcPort, udp.DstPort))
				_, errWrite = writer.WriteString(fmt.Sprintf("UDP Payload:%v\n", strings.TrimFunc(string(udp.Payload), func(r rune) bool {
					return unicode.IsSpace(r)
				})))

			case layers.LayerTypeDNS:
				var resultDNSQuestions, resultDNSAnswers string

				for _, dnsQ := range dns.Questions {
					resultDNSQuestions += string(dnsQ.Name)
				}

				for _, dnsA := range dns.Answers {
					resultDNSAnswers += fmt.Sprintf("%v (%v), %v\n", string(dnsA.Name), dnsA.IP, dnsA.CNAME)
				}

				_, errWrite = writer.WriteString(fmt.Sprintf("DNS questions:'%v', answers:'%v'\n", resultDNSQuestions, resultDNSAnswers))

			case layers.LayerTypeNTP:
				_, errWrite = writer.WriteString(fmt.Sprintf("Version:'%v'\n", ntp.Version))

			case layers.LayerTypeTLS:
				_, errWrite = writer.WriteString(fmt.Sprintf("%v\n", tls.Handshake))

			}

			if errWrite != nil {
				logger.Send("error", fmt.Sprintf("error decode file '%s': %s", fileName, errWrite))
			}
		}
	}

	return err
}

/*
func printApplicationInfo(packet gopacket.Packet) string {
	applicationLayer := packet.ApplicationLayer()
	if applicationLayer != nil {
		return fmt.Sprintf("Application layer/Payload found.\n\t%s\n", applicationLayer.Payload())
	}

	return ""
}

func NetTraffPcapDecoder(filePath, fileName string, fw *os.File, logger commoninterfaces.Logger) error {
	handle, err := pcap.OpenOffline(path.Join(filePath, fileName))
	if err != nil {
		return err
	}

	writer := bufio.NewWriter(fw)
	defer func() {
		if err == nil {
			err = writer.Flush()
		}
	}()

	packets := gopacket.NewPacketSource(handle, layers.LayerTypeEthernet)

	for packet := range packets.Packets() {
		if appStr := printApplicationInfo(packet); appStr != "" {
			writer.WriteString(appStr)

			continue
		}

		writer.WriteString(packet.String())
	}

	return nil
}
*/
