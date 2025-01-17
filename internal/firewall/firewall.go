// Code created by akiidjk on 2025.01.17
package firewall

import (
	"github.com/akiidjk/styx/internal/ebpf"
	"github.com/akiidjk/styx/internal/utils"
	l "github.com/akiidjk/styx/internal/utils/logger"
	"github.com/spf13/cobra"
)

var logger = l.GetLogger()

func init() {

}

var objects = []string{"Dispacher"}

func Run(cmd *cobra.Command, args []string) {
	logger.Info().Msg("Loading eBPF objects")

	ifname, err := cmd.Flags().GetString("ifname")
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to get ifname")
	}

	ebpf.LoadObjectsEBPF(objects)
	defer ebpf.CloseAllObject()

	blockIps, err := cmd.Flags().GetStringSlice("block-ips")
	go ebpf.Dispach(ifname, blockIps)

	utils.HandleTerminate()
}
