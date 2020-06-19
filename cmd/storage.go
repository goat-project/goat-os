package cmd

import (
	"github.com/goat-project/goat-os/constants"
	"github.com/goat-project/goat-os/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var storageFlags = []string{constants.CfgSite}

var storageCmd = &cobra.Command{
	Use:   "storage",
	Short: "Extract storage data",
	Long: "The accounting client is a command-line tool that connects to a cloud, " +
		"extracts data about storages, filters them accordingly and " +
		"then sends them to a server for further processing.",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Init()

		// TODO check if required constants exists
		// TODO set rate limiters
		// TODO account storage
	},
}

func initStorage() {
	goatOsCmd.AddCommand(storageCmd)

	storageCmd.PersistentFlags().String(parseFlagName(constants.CfgSite),
		viper.GetString(constants.CfgSite), "site [SITE]")

	bindFlags(*storageCmd, storageFlags)
}
