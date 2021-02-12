package cmd

import (
	"sync"
	"time"

	"github.com/goat-project/goat-os/auth"
	"github.com/goat-project/goat-os/client"
	"github.com/goat-project/goat-os/constants"
	"github.com/goat-project/goat-os/filter"
	"github.com/goat-project/goat-os/logger"
	"github.com/goat-project/goat-os/preparer"
	"github.com/goat-project/goat-os/processor"
	"github.com/goat-project/goat-os/reader"
	"github.com/goat-project/goat-os/resource/storage"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/time/rate"
)

var storageFlags = []string{constants.CfgSite, constants.CfgAccounted}

var storageRequired []string

var storageDescription = map[string]string{
	constants.CfgSite:      "site [SITE]",
	constants.CfgAccounted: "accounted [storages]",
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

		accounted := viper.GetStringSlice(constants.CfgAccounted)
		if len(accounted) < 1 {
			log.WithField("error", "no accounted for storage are set in configuration").Error("error account storage")
			return
		}

		writeLimiter := rate.NewLimiter(rate.Every(time.Second/time.Duration(requestsPerSecond)), requestsPerSecond)

		var wg sync.WaitGroup

		wg.Add(1)
		go accountStorage(writeLimiter, &wg)
		wg.Wait()
	},
}

func initStorage() {
	goatOsCmd.AddCommand(storageCmd)

	createFlags(storageCmd, storageFlags, storageDescription, storageShorthand)
	bindFlags(*storageCmd, storageFlags)
}

func accountStorage(writeLimiter *rate.Limiter, wg *sync.WaitGroup) {
	defer wg.Done()

	opts := options()

	osClient, err := auth.OpenstackClient(opts)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Fatal("unable to create Openstack client")
	}

	identityClient, err := auth.CreateIdentityV3ServiceClient(osClient)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Fatal("unable to create Identity V3 service client")
	}

	prep := preparer.CreatePreparer(storage.CreatePreparer(reader.CreateReader(identityClient), writeLimiter,
		goatServerConnection()))
	proc := processor.CreateProcessor(storage.CreateProcessor(reader.CreateReader(identityClient)))
	filt := filter.CreateFilter(storage.CreateFilter())

	c := client.Client{}
	c.Run(proc, filt, prep, opts)
}
