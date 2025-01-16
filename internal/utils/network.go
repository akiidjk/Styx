package utils

import (
	"net"
	"os"

	l "github.com/akiidjk/styx/internal/utils/logger"
	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
)

var logger = l.GetLogger()

func LinkInterface(ifname string, program *ebpf.Program) link.Link {
	iface, err := net.InterfaceByName(ifname)
	if err != nil {
		logger.Fatal().Err(err).Str("Ifname", ifname).Msg("Getting interface...")
		os.Exit(1)
	}

	link, err := link.AttachXDP(link.XDPOptions{
		Program:   program,
		Interface: iface.Index,
	})

	if err != nil {
		logger.Fatal().Err(err).Msg("Attaching XDP")
		os.Exit(1)
	}

	return link
}
