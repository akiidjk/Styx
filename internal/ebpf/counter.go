package ebpf

import (
	"os"
	"os/signal"
	"time"

	ebpfGenerated "github.com/akiidjk/styx/internal/ebpf/generated"
	"github.com/akiidjk/styx/internal/utils"
)

func Count() {
	// Load the compiled eBPF ELF and load it into the kernel.
	var objs ebpfGenerated.CounterObjects
	if err := ebpfGenerated.LoadCounterObjects(&objs, nil); err != nil {
		logger.Fatal().Err(err).Msg("Loading eBPF objects")
		os.Exit(1)
	}
	defer objs.Close()

	ifname := "enp5s0" // Change this to an interface on your machine.
	link := utils.LinkInterface(ifname, objs.CountPackets)
	defer link.Close()

	logger.Info().Str("Interface name", ifname).Msg("Counting incoming packets on interface")

	tick := time.Tick(time.Second)
	stop := make(chan os.Signal, 5)
	signal.Notify(stop, os.Interrupt)
	for {
		select {
		case <-tick:
			var count uint64
			err := objs.PktCount.Lookup(uint32(0), &count)
			if err != nil {
				logger.Fatal().Err(err).Msg("Map lookup failed")
				os.Exit(1)
			}
			logger.Info().Uint64("packets", count).Msg("Counter of current packets")
		case <-stop:
			logger.Info().Msg("Received signal, exiting..")
			return
		}
	}
}
