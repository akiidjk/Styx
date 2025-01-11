/*
Copyright © 2024 Akiidjk akiidjk@proton.me
*/
package cmd

import (
	"os"

	ebpfModule "github.com/akiidjk/styx/internal/ebpf"
	"github.com/akiidjk/styx/internal/utils"
	"github.com/spf13/cobra"
)

var ifname string
var blockIPs []string

var rootCmd = &cobra.Command{
	Use:   "styx",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		ebpfModule.RunPacketFilter(ifname, blockIPs)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	go utils.HandleTerminate()

	rootCmd.PersistentFlags().StringVarP(&ifname, "ifname", "i", "", "Interface to be observed")
	rootCmd.MarkFlagRequired("ifname")
	rootCmd.PersistentFlags().StringSliceVarP(&blockIPs, "block-ips", "b", []string{}, "List of IPs to block")
}
