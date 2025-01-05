/*
Copyright © 2024 NAME HERE akiidjk@proton.me
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

	eiud := os.Geteuid()
	if eiud != 0 {
		l.Fatal().Msg("Remember to run the program with `sudo` or with root")
		os.Exit(1)

		l.Debug().Msg("Remember to run the program with `sudo` or with root")
	}

	// Remove resource limits for kernels <5.11.
	if err := rlimit.RemoveMemlock(); err != nil {
		l.Fatal().Err(err).Msg("Removing memlock")
		os.Exit(1)
	}
}

func main() {
	cmd.Execute()
}
