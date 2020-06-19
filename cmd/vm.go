package cmd

import (
	"github.com/goat-project/goat-os/constants"
	"github.com/goat-project/goat-os/logger"
	"github.com/spf13/cobra"
)

var vmFlags = []string{constants.CfgSiteName, constants.CfgCloudType, constants.CfgCloudComputeService}

var vmCmd = &cobra.Command{
	Use:   "vm",
	Short: "Extract virtual machine data",
	Long: "The accounting client is a command-line tool that connects to a cloud, " +
		"extracts data about virtual machines, filters them accordingly and " +
		"then sends them to a server for further processing.",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Init()

		// TODO check if required constants exists
		// TODO set rate limiters
		// TODO account virtual machine
	},
}
