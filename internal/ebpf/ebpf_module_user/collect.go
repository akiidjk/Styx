package ebpfmoduleuser

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"os/signal"

	ebpfModules "github.com/akiidjk/styx/internal/ebpf/generated"
	"github.com/akiidjk/styx/internal/utils"
	"github.com/akiidjk/styx/internal/utils/logger"
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

func processPacket(packetPayload []byte) (uint8, net.IP) {
	var eth layers.Ethernet
	var ip4 layers.IPv4
	var ip6 layers.IPv6
	var tcp layers.TCP

	parser := gopacket.NewDecodingLayerParser(
		layers.LayerTypeEthernet,
		&eth, &ip4, &ip6, &tcp,
	)
	decoded := []gopacket.LayerType{}

	if err := parser.DecodeLayers(packetPayload, &decoded); err != nil {
		logger.Warningf("Could not decode layers: %v", err)
		return 1, ip4.SrcIP
	}

	for _, layerType := range decoded {
		switch layerType {
		case layers.LayerTypeEthernet:
			logger.Debug("Ethernet Layer:")
			logger.Debugf("    Src MAC: %s, Dst MAC: %s", eth.SrcMAC, eth.DstMAC)
		case layers.LayerTypeIPv4:
			logger.Debug("IPv4 Layer:")
			logger.Debugf("    Src IP: %s, Dst IP: %s", ip4.SrcIP, ip4.DstIP)
			if fmt.Sprint(ip4.SrcIP) == "192.168.1.118" {
				return 0, ip4.SrcIP
			}
		case layers.LayerTypeIPv6:
			logger.Debug("IPv6 Layer:")
			logger.Debugf("    Src IP: %s, Dst IP: %s", ip6.SrcIP, ip6.DstIP)
		case layers.LayerTypeTCP:
			logger.Debug("TCP Layer:")
			logger.Debugf("    Src Port: %d, Dst Port: %d", tcp.SrcPort, tcp.DstPort)
		}
	}
	return 1, ip4.SrcIP
}

func Collect(ifname string) {
	var objs ebpfModules.CollecterObjects
	if err := ebpfModules.LoadCollecterObjects(&objs, nil); err != nil {
		logger.Fatal("Loading eBPF objects:", err)
	}
	defer objs.Close()

	link := utils.LinkInterface(ifname, objs.ShareData)
	defer link.Close()

	logger.Info("Sharing incoming packets on ", ifname)

	controlMap := objs.ControlMap
	keyControlMap := uint32(0)
	var valueControlMap uint8 = 255
	var ip net.IP

	stop := make(chan os.Signal, 5)
	signal.Notify(stop, os.Interrupt)

	rd, err := ringbuf.NewReader(objs.PacketMap)
	if err != nil {
		logger.Fatal("Failed to open ringbuf reader: ", err)
		os.Exit(1)
	}
	defer rd.Close()

	go func() {
		for {
			select {
			case <-stop:
				logger.Info("Received signal, exiting..")
				os.Exit(0)
			}
		}
	}()

	for {
		record, err := rd.Read()

		if err != nil {
			logger.Warning("Error reading from ringbuf: ", err)
			continue
		}

		reader := bytes.NewReader(record.RawSample)

		var size int64
		if err := binary.Read(reader, binary.LittleEndian, &size); err != nil {
			logger.Warning("Error reading size: ", err)
			continue
		}

		payload := make([]byte, size)
		if _, err := reader.Read(payload); err != nil {
			logger.Warning("Error reading payload: ", err)
			continue
		}

		packet := Packet{
			Size:    size,
			Payload: payload,
		}

		valueControlMap, ip = processPacket(packet.Payload)
		logger.Debug("Value: ", valueControlMap, " Ip: ", ip)
		if err := controlMap.Update(keyControlMap, valueControlMap, ebpf.UpdateAny); err != nil {
			logger.Fatal("Failed to update control map: ", err)
		}

		valueControlMap = 255 // waiting status
		logger.Debug("Setting waiting: ", valueControlMap)

		// logger.Debug(packetDecoded.Dump())
	}
}
