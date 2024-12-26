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
	fmt.Print(logger.DefaultColors.Magenta + `
███████╗████████╗██╗   ██╗██╗  ██╗    ███████╗██╗██████╗ ███████╗██╗    ██╗ █████╗ ██╗     ██╗
██╔════╝╚══██╔══╝╚██╗ ██╔╝╚██╗██╔╝    ██╔════╝██║██╔══██╗██╔════╝██║    ██║██╔══██╗██║     ██║
███████╗   ██║    ╚████╔╝  ╚███╔╝     █████╗  ██║██████╔╝█████╗  ██║ █╗ ██║███████║██║     ██║
╚════██║   ██║     ╚██╔╝   ██╔██╗     ██╔══╝  ██║██╔══██╗██╔══╝  ██║███╗██║██╔══██║██║     ██║
███████║   ██║      ██║   ██╔╝ ██╗    ██║     ██║██║  ██║███████╗╚███╔███╔╝██║  ██║███████╗███████╗
╚══════╝   ╚═╝      ╚═╝   ╚═╝  ╚═╝    ╚═╝     ╚═╝╚═╝  ╚═╝╚══════╝ ╚══╝╚══╝ ╚═╝  ╚═╝╚══════╝╚══════╝

	created by @akiidjk

` + logger.DefaultColors.Reset)

	eiud := os.Geteuid()
	if eiud != 0 {
		logger.Fatal("Remember to run the program with `sudo` or with root")
		os.Exit(1)
	}

	// Remove resource limits for kernels <5.11.
	if err := rlimit.RemoveMemlock(); err != nil {
		logger.Fatal("Removing memlock:", err)
		os.Exit(1)
	}
}

func main() {
	cmd.Execute()
}
