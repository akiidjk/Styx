package ebpfmoduleuser

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"os/signal"

	ebpfModules "github.com/akiidjk/styx/internal/ebpf/generated"
	"github.com/akiidjk/styx/internal/utils"
	"github.com/akiidjk/styx/internal/utils/logger"
	"github.com/cilium/ebpf/ringbuf"
)

// Set the value in the struct public or the program crash
type packetInfo struct {
	SrcIP    uint32
	DstIP    uint32
	Protocol uint16
	_        uint16
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

	stop := make(chan os.Signal, 5)
	signal.Notify(stop, os.Interrupt)

	rd, err := ringbuf.NewReader(objs.PktData)
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

		var data packetInfo
		if err := binary.Read(bytes.NewReader(record.RawSample), binary.LittleEndian, &data); err != nil {
			logger.Warning("Error decoding data: ", err)
			continue
		}

		logger.Info(fmt.Sprintf("SrcIP: %d.%d.%d.%d, DstIP: %d.%d.%d.%d, Protocol: %d",
			data.SrcIP>>24, (data.SrcIP>>16)&0xFF, (data.SrcIP>>8)&0xFF, data.SrcIP&0xFF,
			data.DstIP>>24, (data.DstIP>>16)&0xFF, (data.DstIP>>8)&0xFF, data.DstIP&0xFF,
			data.Protocol))
	}
}
