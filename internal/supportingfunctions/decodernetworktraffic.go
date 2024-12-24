package supportingfunctions

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcapgo"

	"github.com/av-belyakov/placeholder_ftp/cmd/commoninterfaces"
)

// NetworkTrafficDecoder декодировщик сетевого трафика
func NetworkTrafficDecoder(fileName string, fr, fw *os.File, logger commoninterfaces.Logger) error {
	fmt.Printf("Read file: '%v'\n", fileName)

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
		eth layers.Ethernet
		ip4 layers.IPv4
		ip6 layers.IPv6
		tcp layers.TCP
		udp layers.UDP
		dns layers.DNS
		ntp layers.NTP
		tls layers.TLS
	)

	decoded := []gopacket.LayerType{}
	parser := gopacket.NewDecodingLayerParser(layers.LayerTypeEthernet, &eth, &ip4, &ip6, &tcp, &udp, &dns, &ntp, &tls)

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

		//_, writeErr := writer.WriteString(fmt.Sprintf("%v, packets length: %v\n", ci.Timestamp, ci.CaptureLength))
		//if writeErr != nil {
		//	break
		//}

		e = parser.DecodeLayers(data, &decoded)
		if e != nil {
			continue
		}

		var writeErr error
		for _, layerType := range decoded {
			switch layerType {
			case layers.LayerTypeIPv6:
				_, writeErr = writer.WriteString(fmt.Sprintf("\n%v, packets length: %v\n\tIP6 %v -> %v\n", ci.Timestamp, ci.CaptureLength, ip6.SrcIP, ip6.DstIP))

			case layers.LayerTypeIPv4:
				_, writeErr = writer.WriteString(fmt.Sprintf("\n%v, packets length: %v\n\tIP4 %v -> %v\n", ci.Timestamp, ci.CaptureLength, ip4.SrcIP, ip4.DstIP))

			case layers.LayerTypeTCP:
				payloadSize := len(tcp.Payload)

				_, writeErr = writer.WriteString(fmt.Sprintf("\tTCP port %v -> %v (payload size:'%d')\n", tcp.SrcPort, tcp.DstPort, payloadSize))
				if writeErr != nil {
					continue
				}

				fin := boolToInt8(tcp.FIN)
				syn := boolToInt8(tcp.SYN)
				rst := boolToInt8(tcp.RST)
				psh := boolToInt8(tcp.PSH)
				ack := boolToInt8(tcp.ACK)
				urg := boolToInt8(tcp.URG)

				_, writeErr = writer.WriteString(fmt.Sprintf("\tFlags (FIN:'%v' SYN:'%v' RST:'%v' PSH:'%v' ACK:'%v' URG:'%v')\n", fin, syn, rst, psh, ack, urg))
				if writeErr != nil {
					continue
				}

				if payloadSize != 0 {
					reader := bufio.NewReader(bytes.NewReader(tcp.Payload))

					httpReq, errHTTP := http.ReadRequest(reader)
					if errHTTP == nil {
						proto := httpReq.Proto
						method := httpReq.Method
						//url := httpReq.URL //содержит целый тип, не только значение httpReq.RequestURI но и методы для парсинга запроса
						host := httpReq.Host
						reqURI := httpReq.RequestURI
						userAgent := httpReq.Header.Get("User-Agent")
						//_, writeErr = writer.WriteString(fmt.Sprintf("%v\n", httpReq.Header))

						_, writeErr = writer.WriteString(fmt.Sprintf("\t%v %v %v\n	Host:%v\n	User-Agent:%v\n", proto, method, reqURI, host, userAgent))
						if writeErr != nil {
							continue
						}
					}

					httpRes, errHTTP := http.ReadResponse(reader, httpReq)
					if errHTTP == nil {
						_, writeErr = writer.WriteString(fmt.Sprintf("\tStatus code:%v\n", httpRes.Status))
						if writeErr != nil {
							continue
						}
					}
				}
			case layers.LayerTypeUDP:
				_, writeErr = writer.WriteString(fmt.Sprintf("\tUDP port:%v -> %v\n", udp.SrcPort, udp.DstPort))
			case layers.LayerTypeDNS:
				var resultDNSQuestions, resultDNSAnswers string

				for _, dnsQ := range dns.Questions {
					resultDNSQuestions += string(dnsQ.Name)
				}

				for _, dnsA := range dns.Answers {
					resultDNSAnswers += fmt.Sprintf("%v (%v), %v\n", string(dnsA.Name), dnsA.IP, dnsA.CNAME)
				}

				_, writeErr = writer.WriteString(fmt.Sprintf("\tDNS questions:'%v', answers:'%v'\n", resultDNSQuestions, resultDNSAnswers))
				if writeErr != nil {
					continue
				}
				//_, err = writer.WriteString(fmt.Sprintf("    Questions:'%v', Answers:'%v'\n", dns.Questions, dns.Answers))
			case layers.LayerTypeNTP:
				_, writeErr = writer.WriteString(fmt.Sprintf("\tVersion:'%v'\n", ntp.Version))
			case layers.LayerTypeTLS:
				_, writeErr = writer.WriteString(fmt.Sprintf("\t%v\n", tls.Handshake))

			}

			if writeErr != nil {
				logger.Send("error", fmt.Sprintf("error decode file '%s': %s", fileName, writeErr))
			}
		}
	}

	return err
}
