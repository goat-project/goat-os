package cmd

import (
	"github.com/goat-project/goat-os/constants"
	"github.com/goat-project/goat-os/logger"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var vmFlags = []string{constants.CfgSiteName, constants.CfgCloudType, constants.CfgCloudComputeService}

var vmRequired = []string{constants.CfgSiteName, constants.CfgCloudType}

var vmDescription = map[string]string{
	constants.CfgSiteName:            "site name [VM_SITE_NAME] (required)",
	constants.CfgCloudType:           "cloud type [VM_CLOUD_TYPE] (required)",
	constants.CfgCloudComputeService: "cloud compute service [VM_CLOUD_COMPUTE_SERVICE]",
}

var vmShorthand = map[string]string{}

var vmCmd = &cobra.Command{
	Use:   "vm",
	Short: "Extract virtual machine data",
	Long: "The accounting client is a command-line tool that connects to a cloud, " +
		"extracts data about virtual machines, filters them accordingly and " +
		"then sends them to a server for further processing.",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Init()

		if viper.GetBool("debug") {
			log.WithFields(log.Fields{"version": version}).Debug("goat-os version")
			logFlags(vmFlags)
		}

		err := checkRequired(vmRequired)
		if err != nil {
			log.WithFields(log.Fields{"flag": err}).Fatal("required flag not set")
		}

		// TODO set rate limiters
		// TODO account virtual machine
	},
}

func initVM() {
	goatOsCmd.AddCommand(vmCmd)

	createFlags(vmCmd, vmFlags, vmDescription, vmShorthand)
	bindFlags(*vmCmd, vmFlags)
}
