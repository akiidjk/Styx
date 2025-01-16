/*
Copyright © 2024 akiidjk akiidjk@proton.me
*/
package main

import (
	"fmt"
	"os"

	"github.com/akiidjk/styx/cmd"
	"github.com/akiidjk/styx/internal/utils/logger"
	"github.com/cilium/ebpf/rlimit"
)

func init() {
	logger.SetLevel(0)
	l := logger.GetLogger()
	fmt.Print("\033[35m" + `
███████╗████████╗██╗   ██╗██╗  ██╗    ███████╗██╗██████╗ ███████╗██╗    ██╗ █████╗ ██╗     ██╗
██╔════╝╚══██╔══╝╚██╗ ██╔╝╚██╗██╔╝    ██╔════╝██║██╔══██╗██╔════╝██║    ██║██╔══██╗██║     ██║
███████╗   ██║    ╚████╔╝  ╚███╔╝     █████╗  ██║██████╔╝█████╗  ██║ █╗ ██║███████║██║     ██║
╚════██║   ██║     ╚██╔╝   ██╔██╗     ██╔══╝  ██║██╔══██╗██╔══╝  ██║███╗██║██╔══██║██║     ██║
███████║   ██║      ██║   ██╔╝ ██╗    ██║     ██║██║  ██║███████╗╚███╔███╔╝██║  ██║███████╗███████╗
╚══════╝   ╚═╝      ╚═╝   ╚═╝  ╚═╝    ╚═╝     ╚═╝╚═╝  ╚═╝╚══════╝ ╚══╝╚══╝ ╚═╝  ╚═╝╚══════╝╚══════╝

	created by @akiidjk

` + "\033[0m")

	if os.Geteuid() != 0 {
		l.Fatal().Msg("Run with sudo or as root")
	}

	if err := rlimit.RemoveMemlock(); err != nil {
		l.Fatal().Err(err).Msg("Removing memlock")
	}

	if err := os.Remove("/sys/fs/bpf/shared_events"); err != nil && !os.IsNotExist(err) {
		l.Fatal().Err(err).Msg("Removing shared_events")
	}
}

func main() {
	cmd.Execute()
}
