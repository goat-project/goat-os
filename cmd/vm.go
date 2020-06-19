package cmd

import (
	"github.com/goat-project/goat-os/constants"
	"github.com/goat-project/goat-os/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

func initVM() {
	goatOsCmd.AddCommand(vmCmd)

	vmCmd.PersistentFlags().String(parseFlagName(constants.CfgSiteName), viper.GetString(constants.CfgSiteName),
		"site name [VM_SITE_NAME] (required)")
	vmCmd.PersistentFlags().String(parseFlagName(constants.CfgCloudType), viper.GetString(constants.CfgCloudType),
		"cloud type [VM_CLOUD_TYPE] (required)")
	vmCmd.PersistentFlags().String(parseFlagName(constants.CfgCloudComputeService),
		viper.GetString(constants.CfgCloudComputeService), "cloud compute service [VM_CLOUD_COMPUTE_SERVICE]")

	bindFlags(*vmCmd, vmFlags)
}
