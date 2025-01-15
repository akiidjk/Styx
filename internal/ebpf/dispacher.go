package ebpf

import (
	ebpfGenerated "github.com/akiidjk/styx/internal/ebpf/generated"
	"github.com/akiidjk/styx/internal/utils"
	"github.com/cilium/ebpf"
)

func LoadEBPFObjectsDispacher() (*ebpfGenerated.DispacherObjects, error) {
	var objs ebpfGenerated.DispacherObjects
	if err := ebpfGenerated.LoadDispacherObjects(&objs, nil); err != nil {
		return nil, err
	}
	return &objs, nil
}

func Dispach(ifname string) {

	objDispacher, err := LoadEBPFObjectsDispacher()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to load eBPF objects")
	}
	defer objDispacher.Close()

	ebpfFilterIp, err := LoadEBPFObjectsFilterIp()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to load eBPF objects")
	}
	defer ebpfFilterIp.Close()

	ebpfFilterMac, err := loadEBPFObjectsFilterMac()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to load eBPF objects")
	}
	defer ebpfFilterMac.Close()

	if err = objDispacher.ProgArray.Update(uint32(0), ebpfFilterIp.XdpFilterIp, ebpf.UpdateAny); err != nil {
		logger.Fatal().Err(err).Msg("Failed to update program array")
	}

	if err = objDispacher.ProgArray.Update(uint32(1), ebpfFilterMac.XdpFilterMac, ebpf.UpdateAny); err != nil {
		logger.Fatal().Err(err).Msg("Failed to update program array")
	}

	link := utils.LinkInterface(ifname, objDispacher.XdpDispach)

	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to attach eBPF program")
	}
	defer link.Close()

	for {
	}

}
