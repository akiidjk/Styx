package ebpf

import (
	"bytes"
	"encoding/binary"
	"os"
	"strings"

	ebpfGenerated "github.com/akiidjk/styx/internal/ebpf/generated"
	"github.com/akiidjk/styx/internal/utils"
	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/perf"
)

type LogEvent struct {
	Timestamp uint64
	SrcIp     uint32
	FilterIp  uint32
	Message   [128]byte
}

func LoadEBPFObjects() (*ebpfGenerated.FilteripObjects, error) {
	var objs ebpfGenerated.FilteripObjects
	if err := ebpfGenerated.LoadFilteripObjects(&objs, nil); err != nil {
		return nil, err
	}
	return &objs, nil
}

func AttachEBPFProgram(ifname string, filterProg *ebpf.Program) (link.Link, error) {
	link := utils.LinkInterface(ifname, filterProg)
	return link, nil
}

func SetupPerfReader(perfMap *ebpf.Map) (*perf.Reader, error) {
	reader, err := perf.NewReader(perfMap, os.Getpagesize()*8)
	if err != nil {
		return nil, err
	}
	return reader, nil
}

func UpdateIPFilterMap(ips []string, arrayMap *ebpf.Map) error {
	for i, ipStr := range ips {
		ip, err := utils.IpToDecimal(ipStr)
		if err != nil {
			return err
		}
		key, value := uint32(i), uint32(ip)
		if err := arrayMap.Update(key, value, ebpf.UpdateAny); err != nil {
			return err
		}
	}
	return nil
}

func HandlePerfEvents(reader *perf.Reader) {
	for {
		record, err := reader.Read()
		if err != nil {
			logger.Err(err).Msg("Error reading perf event")
			continue
		}

		var event LogEvent
		if err := binary.Read(bytes.NewReader(record.RawSample), binary.LittleEndian, &event); err != nil {
			logger.Err(err).Msg("Error decoding event")
			continue
		}

		message := strings.TrimSpace(string(event.Message[:]))
		logger.Info().
			Uint64("Timestamp", event.Timestamp).
			Uint32("Src IP", event.SrcIp).
			Uint32("Filter IP", event.FilterIp).
			Str("Message", message).
			Msg("Event received")

	}
}

func RunPacketFilter(ifname string, blockedIPs []string) {
	ebpfObjects, err := LoadEBPFObjects()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to load eBPF objects")
	}
	defer ebpfObjects.Close()

	link, err := AttachEBPFProgram(ifname, ebpfObjects.XdpFilterIp)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to attach eBPF program")
	}
	defer link.Close()

	perfReader, err := SetupPerfReader(ebpfObjects.FilteripMaps.EventOutputMap)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create perf event reader")
	}
	defer perfReader.Close()

	if err := UpdateIPFilterMap(blockedIPs, ebpfObjects.FilteripMaps.IpFilterMap); err != nil {
		logger.Fatal().Err(err).Msg("Failed to update filter map")
	}

	logger.Info().Str("Interface", ifname).Msg("Starting packet filter")
	HandlePerfEvents(perfReader)
}
