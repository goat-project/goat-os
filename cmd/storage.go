package cmd

import (
	"github.com/goat-project/goat-os/constants"
	"github.com/goat-project/goat-os/logger"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var storageFlags = []string{constants.CfgSite}

var storageRequired []string

var storageDescription = map[string]string{
	constants.CfgSite: "site [SITE]",
}

var storageShorthand = map[string]string{}

var storageCmd = &cobra.Command{
	Use:   "storage",
	Short: "Extract storage data",
	Long: "The accounting client is a command-line tool that connects to a cloud, " +
		"extracts data about storages, filters them accordingly and " +
		"then sends them to a server for further processing.",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Init()

		if viper.GetBool("debug") {
			log.WithFields(log.Fields{"version": version}).Debug("goat-os version")
			logFlags(storageFlags)
		}

		err := checkRequired(storageRequired)
		if err != nil {
			log.WithFields(log.Fields{"flag": err}).Fatal("required flag not set")
		}

		// TODO set rate limiters
		// TODO account storage
	},
}

func initStorage() {
	goatOsCmd.AddCommand(storageCmd)

	createFlags(storageCmd, storageFlags, storageDescription, storageShorthand)
	bindFlags(*storageCmd, storageFlags)
}
