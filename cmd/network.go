package cmd

import (
	"github.com/goat-project/goat-os/constants"
	"github.com/goat-project/goat-os/logger"
	"github.com/spf13/cobra"
)

var networkFlags = []string{constants.CfgNetworkSiteName, constants.CfgNetworkCloudType,
	constants.CfgNetworkCloudComputeService}

var networkCmd = &cobra.Command{
	Use:   "network",
	Short: "Extract network data",
	Long: "The accounting client is a command-line tool that connects to a cloud, " +
		"extracts data about networks, filters them accordingly and " +
		"then sends them to a server for further processing.",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Init()

		// TODO check if required constants exists
		// TODO set rate limiters
		// TODO account network
	},
}
