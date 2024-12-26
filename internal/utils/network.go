package utils

import (
	"net"
	"os"

	"github.com/akiidjk/styx/internal/utils/logger"
	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
)

func LinkInterface(ifname string, program *ebpf.Program) link.Link {
	iface, err := net.InterfaceByName(ifname)
	if err != nil {
		logger.Fatal("Getting interface ", ifname, ": ", err)
		os.Exit(1)
	}

	link, err := link.AttachXDP(link.XDPOptions{
		Program:   program,
		Interface: iface.Index,
	})
	if err != nil {
		logger.Fatal("Attaching XDP:", err)
		os.Exit(1)
	}

	return link
}
