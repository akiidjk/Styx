package ebpf

import (
	"bytes"
	"encoding/binary"
	"os"
	"strings"

	ebpfGenerated "github.com/akiidjk/styx/internal/ebpf/generated"
	"github.com/akiidjk/styx/internal/utils"
	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/perf"
)

type LogEvent struct {
	Timestamp uint64
	SrcIp     uint32
	FilterIp  uint32
	Message   [128]byte
}

func LoadEBPFObjectsFilterIp() *ebpfGenerated.FilteripObjects {
	var objs ebpfGenerated.FilteripObjects
	if err := ebpfGenerated.LoadFilteripObjects(&objs, nil); err != nil {
		logger.Fatal().Err(err).Msg("Failed to load eBPF objects")
	}
	return &objs
}

func setupPerfReader(perfMap *ebpf.Map) (*perf.Reader, error) {
	reader, err := perf.NewReader(perfMap, os.Getpagesize()*8)
	if err != nil {
		return nil, err
	}
	return reader, nil
}

func updateIPFilterMap(ips []string, arrayMap *ebpf.Map) error {
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

func handlePerfEvents(reader *perf.Reader) {
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

func RunPacketFilterIP(ifname string, blockedIPs []string) {
	ebpfObjects := LoadEBPFObjectsFilterIp()
	defer ebpfObjects.Close()

	link := utils.LinkInterface(ifname, ebpfObjects.XdpFilterIp)
	defer link.Close()

	perfReader, err := setupPerfReader(ebpfObjects.FilteripMaps.EventOutputMap)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create perf event reader")
	}
	defer perfReader.Close()

	if err := updateIPFilterMap(blockedIPs, ebpfObjects.FilteripMaps.IpFilterMap); err != nil {
		logger.Fatal().Err(err).Msg("Failed to update filter map")
	}

	logger.Info().Str("Interface", ifname).Msg("Starting packet filter")
	handlePerfEvents(perfReader)
}
