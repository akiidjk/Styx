package ebpf

import (
	ebpfGenerated "github.com/akiidjk/styx/internal/ebpf/generated"
)

func loadEBPFObjectsFilterMac() *ebpfGenerated.FiltermacObjects {
	var objs ebpfGenerated.FiltermacObjects
	if err := ebpfGenerated.LoadFiltermacObjects(&objs, nil); err != nil {
		logger.Fatal().Err(err).Msg("Failed to load eBPF objects")
	}
	return &objs
}

func RunPacketFilterMac(ifname string, blockedMacs []string) {
	ebpfObjects := loadEBPFObjectsFilterMac()
	defer ebpfObjects.Close()

	logger.Info().Str("Interface", ifname).Msg("Starting packet filter")
	for {
		// handlePerfEvents(perfReader)
	}
}
