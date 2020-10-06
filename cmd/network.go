package cmd

import (
	"github.com/goat-project/goat-os/constants"
	"github.com/goat-project/goat-os/logger"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var networkFlags = []string{constants.CfgNetworkSiteName, constants.CfgNetworkCloudType,
	constants.CfgNetworkCloudComputeService}

var networkRequired = []string{constants.CfgNetworkSiteName, constants.CfgNetworkCloudType}

var networkDescription = map[string]string{
	constants.CfgNetworkSiteName:            "site name [NETWORK_SITE_NAME] (required)",
	constants.CfgNetworkCloudType:           "cloud type [NETWORK_CLOUD_TYPE] (required)",
	constants.CfgNetworkCloudComputeService: "cloud compute service [NETWORK_CLOUD_COMPUTE_SERVICE]",
}

var networkShorthand = map[string]string{}

var networkCmd = &cobra.Command{
	Use:   "network",
	Short: "Extract network data",
	Long: "The accounting client is a command-line tool that connects to a cloud, " +
		"extracts data about networks, filters them accordingly and " +
		"then sends them to a server for further processing.",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Init()

		if viper.GetBool("debug") {
			log.WithFields(log.Fields{"version": version}).Debug("goat-os version")
			logFlags(networkFlags)
		}

		err := checkRequired(networkRequired)
		if err != nil {
			log.WithFields(log.Fields{"flag": err}).Fatal("required flag not set")
		}

		// TODO set rate limiters
		// TODO account network
	},
}

func initNetwork() {
	goatOsCmd.AddCommand(networkCmd)

	createFlags(networkCmd, networkFlags, networkDescription, networkShorthand)
	bindFlags(*networkCmd, networkFlags)
}
