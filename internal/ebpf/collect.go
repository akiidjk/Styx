package ebpf

import (
	"bytes"
	"encoding/binary"
	"os"
	"os/signal"

	ebpfGenerated "github.com/akiidjk/styx/internal/ebpf/generated"
	"github.com/akiidjk/styx/internal/utils"
	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/ringbuf"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// Set the value in the struct public or the program crash
type Packet struct {
	Payload []byte
	Size    int64
}

func processPacket(packetPayload []byte, ipToBlock string) (uint8, string) {
	var eth layers.Ethernet
	var ip4 layers.IPv4
	var ip6 layers.IPv6
	var tcp layers.TCP
	var udp layers.UDP
	var icmp4 layers.ICMPv4

	parser := gopacket.NewDecodingLayerParser(
		layers.LayerTypeEthernet,
		&eth, &ip4, &ip6, &tcp, &udp, &icmp4,
	)
	decoded := []gopacket.LayerType{}

	if err := parser.DecodeLayers(packetPayload, &decoded); err != nil {
		logger.Warn().Err(err).Msg("Could not decode layers")
		return 1, ""
	}

	for _, layerType := range decoded {
		switch layerType {
		case layers.LayerTypeEthernet:
			logger.Debug().Msg("Ethernet Layer:")
			logger.Debug().Str("Src MAC", eth.SrcMAC.String()).Str("Dest MAC", eth.DstMAC.String())
		case layers.LayerTypeIPv4:
			logger.Debug().Msg("IPv4 Layer:")
			logger.Debug().Str("Src IPv4", ip4.SrcIP.String()).Str("Dest IPv4", ip4.DstIP.String())
			if ip4.SrcIP.String() == ipToBlock {
				return 0, ip4.SrcIP.String()
			}
		case layers.LayerTypeIPv6:
			logger.Debug().Msg("IPv6 Layer:")
			logger.Debug().Str("Src IPv6", ip6.SrcIP.String()).Str("Dest IPv6", ip6.DstIP.String())
		case layers.LayerTypeTCP:
			logger.Debug().Msg("TCP Layer:")
			logger.Debug().Str("TCP Src port", tcp.SrcPort.String()).Str("TCP Dest port", tcp.DstPort.String())
		case layers.LayerTypeUDP:
			logger.Debug().Msg("UDP Layer:")
			logger.Debug().Str("UDP Src port", udp.SrcPort.String()).Str("UDP Dest port", udp.DstPort.String())
		case layers.LayerTypeICMPv4:
			logger.Debug().Msg("ICMPv4 Layer:")
			logger.Debug().Str("TypeCode", icmp4.TypeCode.String()).Uint16("Checksum", icmp4.Checksum)
			if ip4.SrcIP.String() == ipToBlock {
				return 0, ip4.SrcIP.String()
			}
		}
	}

	// Nessun pacchetto corrispondente trovato
	return 1, "2.2.2.2"
}

func Collect(ifname string, ipToBlock string) {
	var objs ebpfGenerated.CollecterObjects
	if err := ebpfGenerated.LoadCollecterObjects(&objs, nil); err != nil {
		logger.Fatal().Err(err).Msg("Loading eBPF objects")
	}
	defer objs.Close()

	link := utils.LinkInterface(ifname, objs.ShareData)
	defer link.Close()

	logger.Info().Str("Interface name", ifname).Msg("Sharing incoming packets on interface")

	controlMap := objs.ControlMap
	keyControlMap := uint32(0)
	var valueControlMap uint8 = 2
	var ip string

	stop := make(chan os.Signal, 5)
	signal.Notify(stop, os.Interrupt)

	rd, err := ringbuf.NewReader(objs.PacketMap)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to open ringbuf reader")
		os.Exit(1)
	}
	defer rd.Close()

	go func() {
		for {
			select {
			case <-stop:
				logger.Info().Msg("Received signal, exiting..")
				os.Exit(0)
			}
		}
	}()

	for {
		record, err := rd.Read()

		if err != nil {
			logger.Warn().Err(err).Msg("Error reading from ringbuf")
			continue
		}

		reader := bytes.NewReader(record.RawSample)

		var size int64
		if err := binary.Read(reader, binary.LittleEndian, &size); err != nil {
			logger.Warn().Err(err).Msg("Error reading size")
			continue
		}

		payload := make([]byte, size)
		if _, err := reader.Read(payload); err != nil {
			logger.Warn().Err(err).Msg("Error reading payload")
			continue
		}

		packet := Packet{
			Size:    size,
			Payload: payload,
		}

		valueControlMap, ip = processPacket(packet.Payload, ipToBlock)
		logger.Debug().Uint8("Value: ", valueControlMap).Str(" Ip: ", ip)
		if err := controlMap.Update(keyControlMap, valueControlMap, ebpf.UpdateAny); err != nil {
			logger.Fatal().Err(err).Msg("Failed to update control map")
			continue
		}

		valueControlMap = 2 // waiting status

		// logger.Debug().Msg(packetDecoded.Dump())
	}
}