package ebpfmoduleuser

import (
	"os"
	"os/signal"
	"time"

	ebpfModules "github.com/akiidjk/fw-ngfw/internal/ebpf/generated"
	"github.com/akiidjk/fw-ngfw/internal/utils"
	"github.com/akiidjk/fw-ngfw/internal/utils/logger"
)

func Count() {
	// Load the compiled eBPF ELF and load it into the kernel.
	var objs ebpfModules.CounterObjects
	if err := ebpfModules.LoadCounterObjects(&objs, nil); err != nil {
		logger.Error("Loading eBPF objects:", err)
	}
	defer objs.Close()

	ifname := "enp5s0" // Change this to an interface on your machine.
	link := utils.LinkInterface(ifname, objs.CountPackets)
	defer link.Close()

	logger.Info("Counting incoming packets on", ifname)

	tick := time.Tick(time.Second)
	stop := make(chan os.Signal, 5)
	signal.Notify(stop, os.Interrupt)
	for {
		select {
		case <-tick:
			var count uint64
			err := objs.PktCount.Lookup(uint32(0), &count)
			if err != nil {
				logger.Error("Map lookup:", err)
			}
			logger.Info("Received", count, "packets")
		case <-stop:
			logger.Info("Received signal, exiting..")
			return
		}
	}
}
