package cmd

import (
	"strings"

	"github.com/goat-project/goat-os/constants"
	"github.com/goat-project/goat-os/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

const version = "1.0.0"

var goatOsFlags = []string{constants.CfgIdentifier, constants.CfgRecordsFrom, constants.CfgRecordsTo,
	constants.CfgRecordsForPeriod, constants.CfgGoatEndpoint, constants.CfgOpenstackIdentityEndpoint,
	constants.CfgUsername, constants.CfgUserID, constants.CfgPassword, constants.CfgPasscode, constants.CfgDomainID,
	constants.CfgDomainName, constants.CfgTenantID, constants.CfgTenantName, constants.CfgAllowReauth,
	constants.CfgTokenID, constants.CfgScopeProjectID, constants.CfgScopeProjectID, constants.CfgScopeDomainID,
	constants.CfgScopeDomainName, constants.CfgScopeSystem, constants.CfgAppCredentialID, constants.CfgAppCredentialName,
	constants.CfgAppCredentialSecret, constants.CfgOpenstackTimeout, constants.CfgDebug, constants.CfgLogPath}

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

		// TODO check if required constants from config exists
		// TODO set rate limiters
		// TODO account vm, network, storage
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
}

func initGoatOs() {
	cobra.OnInitialize(initConfig)

	goatOsCmd.PersistentFlags().StringP(constants.CfgIdentifier, "i", viper.GetString(constants.CfgIdentifier),
		"goat identifier [IDENTIFIER] (required)")

	goatOsCmd.PersistentFlags().StringP(constants.CfgRecordsFrom, "f", viper.GetString(constants.CfgRecordsFrom),
		"records from [TIME]")
	goatOsCmd.PersistentFlags().StringP(constants.CfgRecordsTo, "t", viper.GetString(constants.CfgRecordsTo),
		"records to [TIME]")
	goatOsCmd.PersistentFlags().StringP(constants.CfgRecordsForPeriod, "p",
		viper.GetString(constants.CfgRecordsForPeriod), "records for period [TIME PERIOD]")

	goatOsCmd.PersistentFlags().StringP(constants.CfgGoatEndpoint, "e", viper.GetString(constants.CfgGoatEndpoint),
		"goat server [GOAT_SERVER_ENDPOINT] (required)")
	goatOsCmd.PersistentFlags().StringP(constants.CfgOpenstackIdentityEndpoint, "o",
		viper.GetString(constants.CfgOpenstackIdentityEndpoint),
		"Openstack identity endpoint [OS_IDENTITY_ENDPOINT] (required)")

	goatOsCmd.PersistentFlags().String(constants.CfgUsername, viper.GetString(constants.CfgUsername),
		"Openstack authentication username [OS_USERNAME]")
	goatOsCmd.PersistentFlags().String(constants.CfgUserID, viper.GetString(constants.CfgUserID),
		"Openstack authentication user ID [OS_USER_ID]")
	goatOsCmd.PersistentFlags().String(constants.CfgPassword, viper.GetString(constants.CfgPassword),
		"Openstack authentication password [OS_PASSWORD]")
	goatOsCmd.PersistentFlags().String(constants.CfgPasscode, viper.GetString(constants.CfgPasscode),
		"Openstack authentication passcoe [OS_PASSCODE]")
	goatOsCmd.PersistentFlags().String(constants.CfgDomainID, viper.GetString(constants.CfgDomainID),
		"Openstack authentication domain ID [OS_DOMAIN_ID]")
	goatOsCmd.PersistentFlags().String(constants.CfgDomainName, viper.GetString(constants.CfgDomainName),
		"Openstack authentication domain name [OS_DOMAIN_NAME]")
	goatOsCmd.PersistentFlags().String(constants.CfgTenantID, viper.GetString(constants.CfgTenantID),
		"Openstack authentication tenant ID [OS_TENANT_ID]")
	goatOsCmd.PersistentFlags().String(constants.CfgTenantName, viper.GetString(constants.CfgTenantName),
		"Openstack authentication tenant name [OS_TENANT_NAME]")
	goatOsCmd.PersistentFlags().String(constants.CfgAllowReauth, viper.GetString(constants.CfgAllowReauth),
		"Openstack authentication allow reauth. [OS_ALLOW_REAUTH]")
	goatOsCmd.PersistentFlags().String(constants.CfgTokenID, viper.GetString(constants.CfgTokenID),
		"Openstack authentication token ID [OS_TOKEN_ID]")
	goatOsCmd.PersistentFlags().String(constants.CfgScopeProjectID, viper.GetString(constants.CfgScopeProjectID),
		"Openstack scope project ID [OS_SCOPE_PROJECT_ID]")
	goatOsCmd.PersistentFlags().String(constants.CfgScopeProjectName, viper.GetString(constants.CfgScopeProjectName),
		"Openstack scope project name [OS_SCOPE_PROJECT_NAME]")
	goatOsCmd.PersistentFlags().String(constants.CfgScopeDomainID, viper.GetString(constants.CfgScopeDomainID),
		"Openstack scope domain ID [OS_SCOPE_DOMAIN_ID]")
	goatOsCmd.PersistentFlags().String(constants.CfgScopeDomainName, viper.GetString(constants.CfgScopeDomainName),
		"Openstack scope domain name [OS_SCOPE_DOMAIN_NAME]")
	goatOsCmd.PersistentFlags().String(constants.CfgScopeSystem, viper.GetString(constants.CfgScopeSystem),
		"Openstack scope system [OS_SCOPE_SYSTEM]")
	goatOsCmd.PersistentFlags().String(constants.CfgAppCredentialID, viper.GetString(constants.CfgAppCredentialID),
		"Openstack application credential ID [OS_APPCREDENTIAL_ID]")
	goatOsCmd.PersistentFlags().String(constants.CfgAppCredentialName, viper.GetString(constants.CfgAppCredentialName),
		"Openstack application credential name [OS_APPCREDENTIAL_NAME]")
	goatOsCmd.PersistentFlags().String(constants.CfgAppCredentialSecret, viper.GetString(constants.CfgAppCredentialSecret),
		"Openstack application credential secret [OS_APPCREDENTIAL_SECRET]")

	goatOsCmd.PersistentFlags().String(constants.CfgOpenstackTimeout, viper.GetString(constants.CfgOpenstackTimeout),
		"timeout for Openstack calls [TIMEOUT_FOR_OPENSTACK_CALLS] (required)")

	goatOsCmd.PersistentFlags().StringP(constants.CfgDebug, "d", viper.GetString(constants.CfgDebug),
		"debug")
	goatOsCmd.PersistentFlags().String(constants.CfgLogPath, viper.GetString(constants.CfgLogPath), "path to log file")

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

func bindFlags(command cobra.Command, flagsForBinding []string) {
	for _, flag := range flagsForBinding {
		err := viper.BindPFlag(flag, command.PersistentFlags().Lookup(parseFlagName(flag)))
		if err != nil {
			log.WithFields(log.Fields{"error": err, "flag": flag}).Panic("unable to initialize flag")
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
