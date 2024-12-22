package main

import (
	"fmt"
	"log"
	"os"

	ebpfmoduleuser "github.com/akiidjk/fw-ngfw/internal/ebpf/ebpf_module_user"
	"github.com/akiidjk/fw-ngfw/internal/utils"
	"github.com/akiidjk/fw-ngfw/internal/utils/logger"
	"github.com/cilium/ebpf/rlimit"
)

func init() {
	logger.SetLevel(0)
	fmt.Print(utils.Magenta + `
███████╗████████╗██╗   ██╗██╗  ██╗    ███████╗██╗██████╗ ███████╗██╗    ██╗ █████╗ ██╗     ██╗
██╔════╝╚══██╔══╝╚██╗ ██╔╝╚██╗██╔╝    ██╔════╝██║██╔══██╗██╔════╝██║    ██║██╔══██╗██║     ██║
███████╗   ██║    ╚████╔╝  ╚███╔╝     █████╗  ██║██████╔╝█████╗  ██║ █╗ ██║███████║██║     ██║
╚════██║   ██║     ╚██╔╝   ██╔██╗     ██╔══╝  ██║██╔══██╗██╔══╝  ██║███╗██║██╔══██║██║     ██║
███████║   ██║      ██║   ██╔╝ ██╗    ██║     ██║██║  ██║███████╗╚███╔███╔╝██║  ██║███████╗███████╗
╚══════╝   ╚═╝      ╚═╝   ╚═╝  ╚═╝    ╚═╝     ╚═╝╚═╝  ╚═╝╚══════╝ ╚══╝╚══╝ ╚═╝  ╚═╝╚══════╝╚══════╝

	created by @akiidjk

` + utils.Reset)
}

func main() {
	// Remove resource limits for kernels <5.11.
	if err := rlimit.RemoveMemlock(); err != nil {
		log.Fatal("Removing memlock:", err)
	}

	ifname := os.Args[1]
	ebpfmoduleuser.Collect(ifname)
	// ebpfmoduleuser.Count()
}
