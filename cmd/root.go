/*
Copyright © 2024 NAME HERE akiidjk@proton.me
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"

	ebpfmoduleuser "github.com/akiidjk/styx/internal/ebpf/ebpf_module_user"
)

var ifname string
var ipToBlock string

var rootCmd = &cobra.Command{
	Use:   "styx",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		ebpfmoduleuser.Collect(ifname, ipToBlock)
		// ebpfmoduleuser.Count()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Flags defining
	rootCmd.PersistentFlags().StringVarP(&ifname, "ifname", "i", "", "Interface to be observed")
	rootCmd.MarkFlagsOneRequired("ifname")
	rootCmd.PersistentFlags().StringVarP(&ipToBlock, "rule", "r", "", "Ip to block")
}
