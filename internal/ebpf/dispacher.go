package ebpf

import (
	ebpfGenerated "github.com/akiidjk/styx/internal/ebpf/generated"
	"github.com/akiidjk/styx/internal/utils"
	"github.com/cilium/ebpf"
)

const pinPathProgArray string = "/sys/fs/bpf/prog_array"
const pinPathEventOutputMap string = "/sys/fs/bpf/event_output"

var objects = []string{"FilterIp", "FilterMac"}

/*
Load the eBPF objects for the dispacher program
*/
func LoadEBPFObjectsDispacher() *ebpfGenerated.DispacherObjects {
	var objs ebpfGenerated.DispacherObjects
	opts := ebpf.CollectionOptions{
		Maps: ebpf.MapOptions{
			PinPath: "/sys/fs/bpf",
		},
	}
	if err := ebpfGenerated.LoadDispacherObjects(&objs, &opts); err != nil {
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

func Dispach(ifname string, blockIps []string) {
	logger.Info().Msgf("Starting Dispach function for interface: %s", ifname)
	dispacherObj, err := GetObject[*ebpfGenerated.DispacherObjects]("Dispacher")
	if err != nil {
		logger.Fatal().Err(err).Msgf("Failed to get Dispacher object for interface %s", ifname)
	}
	logger.Info().Msg("Successfully obtained Dispacher object")

	logger.Info().Msgf("Linking interface %s with XDP program", ifname)
	link := utils.LinkInterface(ifname, (*dispacherObj).XdpDispach)
	defer func() {
		logger.Info().Msg("Closing interface link")
		link.Close()
	}()
	logger.Info().Msg("Interface successfully linked")

	// Load pinned map
	ProgArray, err := ebpf.LoadPinnedMap(pinPathProgArray, nil)
	if err != nil {
		logger.Fatal().Err(err).Msgf("Failed to load pinned map from path: %s", pinPathProgArray)
	}
	logger.Info().Msgf("Successfully loaded pinned map from: %s", pinPathProgArray)
	defer func() {
		logger.Info().Msg("Closing pinned map")
		ProgArray.Close()
	}()

	EventOutputMap, err := ebpf.LoadPinnedMap(pinPathEventOutputMap, nil)
	if err != nil {
		logger.Fatal().Err(err).Msgf("Failed to load pinned map from path: %s", EventOutputMap)
	}
	logger.Info().Msgf("Successfully loaded pinned map from: %s", EventOutputMap)
	defer func() {
		logger.Info().Msg("Closing pinned map")
		EventOutputMap.Close()
	}()

	LoadObjectsEBPF(objects)

	logger.Info().Msg("Adding programs to map...")
	addPrograms(GetPrograms(), ProgArray)
	logger.Info().Msg("Programs successfully added to map")

	logger.Info().Msg("Entering infinite loop to keep program running")

	RunPacketFilterIP(ifname, blockIps)
	for {
		// time.Sleep(time.Second)
	}
}
