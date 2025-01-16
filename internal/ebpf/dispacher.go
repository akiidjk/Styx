package ebpf

import (
	ebpfGenerated "github.com/akiidjk/styx/internal/ebpf/generated"
	"github.com/akiidjk/styx/internal/utils"
	"github.com/cilium/ebpf"
)

/*
Load the eBPF objects for the dispacher program
*/
func LoadEBPFObjectsDispacher() *ebpfGenerated.DispacherObjects {
	var objs ebpfGenerated.DispacherObjects
	if err := ebpfGenerated.LoadDispacherObjects(&objs, nil); err != nil {
		logger.Fatal().Err(err).Msg("Failed to load eBPF objects")
	}
	return &objs
}

/*
Update the program array with the new programs
*/
func addPrograms(programs []*ebpf.Program, programMap *ebpf.Map) {

	for index, program := range programs {
		if err := programMap.Update(uint32(index), program, ebpf.UpdateAny); err != nil {
			logger.Fatal().Err(err).Msg("Failed to update program array")
		}
	}
}

func Dispach(ifname string) {
	dispacherObj, err := GetObject[*ebpfGenerated.DispacherObjects]("Dispacher")
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to get Dispacher object")
	}
	addPrograms(GetPrograms(), (*dispacherObj).ProgArray)
	link := utils.LinkInterface(ifname, (*dispacherObj).XdpDispach)
	defer link.Close()

	for {
	}
}
