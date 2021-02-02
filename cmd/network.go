package cmd

import (
	"sync"
	"time"

	"github.com/goat-project/goat-os/auth"
	"github.com/goat-project/goat-os/reader"

	"github.com/goat-project/goat-os/client"
	"github.com/goat-project/goat-os/constants"
	"github.com/goat-project/goat-os/filter"
	"github.com/goat-project/goat-os/logger"
	"github.com/goat-project/goat-os/preparer"
	"github.com/goat-project/goat-os/processor"
	"github.com/goat-project/goat-os/resource/network"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/time/rate"
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

		writeLimiter := rate.NewLimiter(rate.Every(time.Second/time.Duration(requestsPerSecond)), requestsPerSecond)

		var wg sync.WaitGroup

		wg.Add(1)
		go accountNetwork(writeLimiter, &wg)
		wg.Wait()
	},
}

func initNetwork() {
	goatOsCmd.AddCommand(networkCmd)

	createFlags(networkCmd, networkFlags, networkDescription, networkShorthand)
	bindFlags(*networkCmd, networkFlags)
}

func accountNetwork(writeLimiter *rate.Limiter, wg *sync.WaitGroup) {
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

	proc := processor.CreateProcessor(network.CreateProcessor(reader.CreateReader(identityClient)))
	filt := filter.CreateFilter(network.CreateFilter())
	prep := preparer.CreatePreparer(network.CreatePreparer(writeLimiter, goatServerConnection()))

	c := client.Client{}

	c.Run(proc, filt, prep, opts)
}
