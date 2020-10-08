package constants

// global constants
const (
	// CfgIdentifier represents string identifier of a goat-os instance
	CfgIdentifier = "identifier"

	// CfgRecordsFrom represents time which records are filtered from
	CfgRecordsFrom = "records-from"
	// CfgRecordsTo represents time which records are filtered to
	CfgRecordsTo = "records-to"
	// CfgRecordsFrom represents duration which records are filtered for
	CfgRecordsForPeriod = "records-for-period"

	// CfgGoatEndpoint represents string of goat server endpoint
	CfgGoatEndpoint = "endpoint"

	// CfgOpenstackIdentityEndpoint represents string of Openstack identity endpoint
	CfgOpenstackIdentityEndpoint = "openstack-identity-endpoint"

	// CfgDebug represents true for debug mode; false otherwise
	CfgDebug = "debug"

	// CfgLogPath represents path to log file
	CfgLogPath = "log-path"
)
