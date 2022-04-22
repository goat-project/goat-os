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
	"github.com/goat-project/goat-os/resource/gpu"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/time/rate"
)

var gpuFlags []string

var gpuRequired = []string{constants.CfgGPUSiteName}

var gpuDescription = map[string]string{
	constants.CfgGPUSiteName: "site name [SITE]",
}

var gpuShorthand = map[string]string{}

var gpuCmd = &cobra.Command{
	Use:   "gpu",
	Short: "Extract gpu data",
	Long: "The accounting client is a command-line tool that connects to a cloud, " +
		"extracts data about gpus, filters them accordingly and " +
		"then sends them to a server for further processing.",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Init()

		if viper.GetBool("debug") {
			log.WithFields(log.Fields{"version": version}).Debug("goat-os version")
			logFlags(gpuFlags)
		}

		err := checkRequired(gpuRequired)
		if err != nil {
			log.WithFields(log.Fields{"flag": err}).Fatal("required flag not set")
		}

		accounted := viper.GetStringSlice(constants.CfgAccounted)
		if len(accounted) < 1 {
			log.WithField("error", "no accounted for gpu are set in configuration").Error("error account gpu")
			return
		}

		writeLimiter := rate.NewLimiter(rate.Every(time.Second/time.Duration(requestsPerSecond)), requestsPerSecond)

		var wg sync.WaitGroup

		wg.Add(1)
		go accountGPU(writeLimiter, &wg)
		wg.Wait()
	},
}

func initGPU() {
	goatOsCmd.AddCommand(gpuCmd)

	createFlags(gpuCmd, gpuFlags, gpuDescription, gpuShorthand)
	bindFlags(*gpuCmd, gpuFlags)
}

func accountGPU(writeLimiter *rate.Limiter, wg *sync.WaitGroup) {
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

	prep := preparer.CreatePreparer(gpu.CreatePreparer(reader.CreateReader(identityClient), reader.CreateReader(computeClient),
		writeLimiter, goatServerConnection()))
	proc := processor.CreateProcessor(gpu.CreateProcessor(reader.CreateReader(identityClient)))
	filt := filter.CreateFilter(gpu.CreateFilter())

	c := client.Client{}
	c.Run(proc, filt, prep, opts)
}
