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
	"github.com/goat-project/goat-os/resource/server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/time/rate"
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

		writeLimiter := rate.NewLimiter(rate.Every(time.Second/time.Duration(requestsPerSecond)), requestsPerSecond)

		var wg sync.WaitGroup

		wg.Add(1)
		go accountVM(writeLimiter, &wg)
		wg.Wait()
	},
}

func initVM() {
	goatOsCmd.AddCommand(vmCmd)

	createFlags(vmCmd, vmFlags, vmDescription, vmShorthand)
	bindFlags(*vmCmd, vmFlags)
}

func accountVM(writeLimiter *rate.Limiter, wg *sync.WaitGroup) {
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

	computeClient, err := auth.CreateComputeV2ServiceClient(osClient)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Fatal("unable to create Compute V2 service client")
	}

	prep := preparer.CreatePreparer(server.CreatePreparer(reader.CreateReader(identityClient),
		reader.CreateReader(computeClient), writeLimiter, goatServerConnection()))
	proc := processor.CreateProcessor(server.CreateProcessor(reader.CreateReader(identityClient)))
	filt := filter.CreateFilter(server.CreateFilter())

	c := client.Client{}
	c.Run(proc, filt, prep, opts)
}
