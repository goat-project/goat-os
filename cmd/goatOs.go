package cmd

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc/credentials/insecure"

	"golang.org/x/time/rate"

	"github.com/gophercloud/gophercloud"

	"google.golang.org/grpc"

	"github.com/goat-project/goat-os/constants"
	"github.com/goat-project/goat-os/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

const version = "1.0.0"
const requestsPerSecond = 30

var goatOsFlags = []string{constants.CfgIdentifier, constants.CfgRecordsFrom, constants.CfgRecordsTo,
	constants.CfgRecordsForPeriod, constants.CfgGoatEndpoint, constants.CfgOpenstackIdentityEndpoint,
	constants.CfgUsername, constants.CfgUserID, constants.CfgPassword,
	constants.CfgPasscode, constants.CfgDomainID, constants.CfgDomainName, constants.CfgTenantID,
	constants.CfgTenantName, constants.CfgAllowReauth, constants.CfgTokenID, constants.CfgScopeProjectID,
	constants.CfgScopeProjectName, constants.CfgScopeDomainID, constants.CfgScopeDomainName, constants.CfgScopeSystem,
	constants.CfgAppCredentialID, constants.CfgAppCredentialName, constants.CfgAppCredentialSecret,
	constants.CfgEndpointType, constants.CfgEndpointName, constants.CfgEndpointRegion,
	constants.CfgEndpointAvailability, constants.CfgDebug, constants.CfgLogPath}

var goatOsRequired = []string{constants.CfgIdentifier, constants.CfgGoatEndpoint,
	constants.CfgOpenstackIdentityEndpoint}

var goatOsDescription = map[string]string{
	constants.CfgIdentifier:       "goat identifier [IDENTIFIER] (required)",
	constants.CfgRecordsFrom:      "records from [TIME]",
	constants.CfgRecordsTo:        "records to [TIME]",
	constants.CfgRecordsForPeriod: "records for period [TIME PERIOD]",

	constants.CfgGoatEndpoint:              "goat server [GOAT_SERVER_ENDPOINT] (required)",
	constants.CfgOpenstackIdentityEndpoint: "Openstack identity endpoint [OS_IDENTITY_ENDPOINT] (required)",

	constants.CfgUsername:            "Openstack authentication username [OS_USERNAME]",
	constants.CfgUserID:              "Openstack authentication user TenantID [OS_USER_ID]",
	constants.CfgPassword:            "Openstack authentication password [OS_PASSWORD]",
	constants.CfgPasscode:            "Openstack authentication passcode [OS_PASSCODE]",
	constants.CfgDomainID:            "Openstack authentication domain TenantID [OS_DOMAIN_ID]",
	constants.CfgDomainName:          "Openstack authentication domain name [OS_DOMAIN_NAME]",
	constants.CfgTenantID:            "Openstack authentication tenant TenantID [OS_TENANT_ID]",
	constants.CfgTenantName:          "Openstack authentication tenant name [OS_TENANT_NAME]",
	constants.CfgAllowReauth:         "Openstack authentication allow reauth. [OS_ALLOW_REAUTH]",
	constants.CfgTokenID:             "Openstack authentication token TenantID [OS_TOKEN_ID]",
	constants.CfgScopeProjectID:      "Openstack scope project TenantID [OS_SCOPE_PROJECT_ID]",
	constants.CfgScopeProjectName:    "Openstack scope project name [OS_SCOPE_PROJECT_NAME]",
	constants.CfgScopeDomainID:       "Openstack scope domain TenantID [OS_SCOPE_DOMAIN_ID]",
	constants.CfgScopeDomainName:     "Openstack scope domain name [OS_SCOPE_DOMAIN_NAME]",
	constants.CfgScopeSystem:         "Openstack scope system [OS_SCOPE_SYSTEM]",
	constants.CfgAppCredentialID:     "Openstack application credential TenantID [OS_APPCREDENTIAL_ID]",
	constants.CfgAppCredentialName:   "Openstack application credential name [OS_APPCREDENTIAL_NAME]",
	constants.CfgAppCredentialSecret: "Openstack application credential secret [OS_APPCREDENTIAL_SECRET]",

	constants.CfgEndpointType:         "Openstack endpoint type [OS_ENDPOINT_TYPE]",
	constants.CfgEndpointName:         "Openstack endpoint name [OS_ENDPOINT_NAME]",
	constants.CfgEndpointRegion:       "Openstack endpoint region [OS_ENDPOINT_REGION]",
	constants.CfgEndpointAvailability: "Openstack endpoint availability (public, internal, admin)",

	constants.CfgDebug:   "debug",
	constants.CfgLogPath: "path to log file",
}

var goatOsShorthand = map[string]string{
	constants.CfgIdentifier:                "i",
	constants.CfgRecordsFrom:               "f",
	constants.CfgRecordsTo:                 "t",
	constants.CfgRecordsForPeriod:          "p",
	constants.CfgGoatEndpoint:              "e",
	constants.CfgOpenstackIdentityEndpoint: "o",
	constants.CfgDebug:                     "d",
}

var goatOsCmd = &cobra.Command{
	Use:   "goat-os",
	Short: "extracts data about virtual machines, networks and storages",
	Long: "The accounting client is a command-line tool that connects to a cloud, " +
		"extracts data about virtual machines, networks and storages, filters them accordingly and " +
		"then sends them to a server for further processing.",
	Version: version,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Init()

		if viper.GetBool("debug") {
			log.WithFields(log.Fields{"version": version}).Debug("goat-os version")
			logFlags(append(vmFlags, append(networkFlags, storageFlags...)...))
		}

		err := checkRequired(append(goatOsRequired, append(vmRequired, append(networkRequired, storageRequired...)...)...))
		if err != nil {
			log.WithFields(log.Fields{"flag": err}).Fatal("required flag not set")
		}

		writeLimiter := rate.NewLimiter(rate.Every(time.Second/time.Duration(requestsPerSecond)), requestsPerSecond)

		var wg sync.WaitGroup

		wg.Add(4)
		go accountVM(writeLimiter, &wg)
		go accountNetwork(writeLimiter, &wg)
		go accountStorage(writeLimiter, &wg)
		go accountGPU(writeLimiter, &wg)
		wg.Wait()
	},
}

// Execute uses the args (os.Args[1:] by default)
// and run through the command tree finding appropriate matches
// for commands and then corresponding flags.
func Execute() {
	if err := goatOsCmd.Execute(); err != nil {
		log.WithFields(log.Fields{"error": err}).Fatal("fatal error execute")
	}
}

// Initialize initializes configuration and CLI options.
func Initialize() {
	initGoatOs()

	initVM()
	initNetwork()
	initStorage()
	initGPU()
}

func initGoatOs() {
	cobra.OnInitialize(initConfig)

	createFlags(goatOsCmd, goatOsFlags, goatOsDescription, goatOsShorthand)
	bindFlags(*goatOsCmd, goatOsFlags)

	viper.SetDefault("author", "Lenka Svetlovska")
	viper.SetDefault("license", "Apache-2.0 License")
}

func initConfig() {
	// name of config file (without extension)
	viper.SetConfigName("goat-os")

	// paths to look for the config file in
	viper.AddConfigPath("config/")
	viper.AddConfigPath("/etc/goat-os/")
	viper.AddConfigPath("$HOME/.goat-os/")

	// find and read the config file
	err := viper.ReadInConfig()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error config file")
	}
}

func createFlags(command *cobra.Command, flags []string, descriptions map[string]string, shorthands map[string]string) {
	for _, flag := range flags {
		if shorthands[flag] != "" {
			command.PersistentFlags().StringP(parseFlagName(flag), shorthands[flag], viper.GetString(flag), descriptions[flag])
		} else {
			command.PersistentFlags().String(parseFlagName(flag), viper.GetString(flag), descriptions[flag])
		}
	}
}

func bindFlags(command cobra.Command, flagsForBinding []string) {
	for _, flag := range flagsForBinding {
		err := viper.BindPFlag(flag, command.PersistentFlags().Lookup(parseFlagName(flag)))
		if err != nil {
			log.WithFields(log.Fields{"error": err, "flag": flag}).Fatal("unable to initialize flag")
		}
	}
}

func parseFlagName(cfgName string) string {
	return lastString(strings.Split(cfgName, "."))
}

func lastString(ss []string) string {
	// This should not happen since it is passing a predefined non-empty strings.
	// It panic here since this will happen only if a mistake in code is made.
	if len(ss) == 0 {
		log.Panic("parsing empty string")
	}

	return ss[len(ss)-1]
}

func logFlags(flags []string) {
	for _, flag := range append(goatOsFlags, flags...) {
		log.WithFields(log.Fields{"flag": flag, "value": viper.Get(flag)}).Debug("flag initialized")
	}
}

func checkRequired(required []string) error {
	for _, req := range required {
		if viper.GetString(req) == "" {
			return fmt.Errorf(req)
		}
	}

	return nil
}

func goatServerConnection() *grpc.ClientConn {
	conn, err := grpc.Dial(viper.GetString(constants.CfgGoatEndpoint), grpc.WithTransportCredentials(
		insecure.NewCredentials()))
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Fatal("error connect to goat server via gRPC")
	}

	return conn
}

func options() gophercloud.AuthOptions {
	return gophercloud.AuthOptions{
		IdentityEndpoint: viper.GetString(constants.CfgOpenstackIdentityEndpoint),
		Username:         viper.GetString(constants.CfgUsername),
		UserID:           viper.GetString(constants.CfgUserID),
		Password:         viper.GetString(constants.CfgPassword),
		Passcode:         viper.GetString(constants.CfgPasscode),
		DomainID:         viper.GetString(constants.CfgDomainID),
		DomainName:       viper.GetString(constants.CfgDomainName),
		TenantID:         viper.GetString(constants.CfgTenantID),
		TenantName:       viper.GetString(constants.CfgTenantName),
		AllowReauth:      viper.GetBool(constants.CfgAllowReauth),
		TokenID:          viper.GetString(constants.CfgTokenID),
		Scope: &gophercloud.AuthScope{
			ProjectID:   viper.GetString(constants.CfgScopeProjectID),
			ProjectName: viper.GetString(constants.CfgScopeProjectName),
			DomainID:    viper.GetString(constants.CfgScopeDomainID),
			DomainName:  viper.GetString(constants.CfgScopeDomainName),
			System:      viper.GetBool(constants.CfgScopeSystem),
		},
		ApplicationCredentialID:     viper.GetString(constants.CfgAppCredentialID),
		ApplicationCredentialName:   viper.GetString(constants.CfgAppCredentialName),
		ApplicationCredentialSecret: viper.GetString(constants.CfgAppCredentialSecret),
	}
}
