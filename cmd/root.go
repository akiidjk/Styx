/*
Copyright © 2024 Akiidjk akiidjk@proton.me
*/
package cmd

import (
	"github.com/akiidjk/styx/internal/firewall"
	"github.com/akiidjk/styx/internal/utils/logger"
	"github.com/spf13/cobra"
)

var ifname string
var blockIPs []string
var blockMacs []string

var rootCmd = &cobra.Command{
	Use:   "styx",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: firewall.Run,
}

var l = logger.GetLogger()

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		l.Fatal().Err(err).Msg("Failed to execute root command")
	}
}

func init() {
	l.Info().Msg("Starting Styx")
	rootCmd.PersistentFlags().StringVarP(&ifname, "ifname", "i", "", "Interface to be observed")
	rootCmd.MarkFlagRequired("ifname")
	rootCmd.PersistentFlags().StringSliceVarP(&blockIPs, "block-ips", "", []string{}, "List of IPs address to block")
	rootCmd.PersistentFlags().StringSliceVarP(&blockMacs, "block-macs", "", []string{}, "List of MACs address to block")
	l.Debug().Msg("Root command initialized")
}
