package ebpf

import (
	"log"
	"time"

	ebpfGenerated "github.com/akiidjk/styx/internal/ebpf/generated"
	"github.com/akiidjk/styx/internal/utils"
)

func loadEBPFObjectsFilterMac() (*ebpfGenerated.FiltermacObjects, error) {
	var objs ebpfGenerated.FiltermacObjects
	if err := ebpfGenerated.LoadFiltermacObjects(&objs, nil); err != nil {
		return nil, err
	}
	return &objs, nil
}

func RunPacketFilterMac(ifname string, blockedMacs []string) {
	ebpfObjects, err := loadEBPFObjectsFilterMac()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to load eBPF objects")
	}
	defer ebpfObjects.Close()

	start := time.Now()
	log.Println("ATTACO")
	link := utils.LinkInterface(ifname, ebpfObjects.XdpFilterMac)
	log.Println("STACCO")
	link.Close()
	elapsed := time.Since(start).Milliseconds()
	log.Printf("Tempo trascorso: %d ms\n", elapsed)

	log.Println("ATACCO")
	link = utils.LinkInterface(ifname, ebpfObjects.XdpFilterMac)
	log.Println("STACCO")
	link.Close()
	elapsed = time.Since(start).Milliseconds()
	log.Printf("Tempo trascorso: %d ms\n", elapsed)

	log.Println("ATTACO")
	link = utils.LinkInterface(ifname, ebpfObjects.XdpFilterMac)
	log.Println("STACCO")
	link.Close()
	elapsed = time.Since(start).Milliseconds()
	log.Printf("Tempo trascorso: %d ms\n", elapsed)

	logger.Info().Str("Interface", ifname).Msg("Starting packet filter")
	for {
		// handlePerfEvents(perfReader)
	}
}
