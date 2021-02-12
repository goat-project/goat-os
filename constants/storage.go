package constants

// prefix for storage subcommands
const cfgStoragePrefix = "storage."

// constants for storage subcommand
const (
	// CfgSite represents string of storage site
	CfgSite = cfgStoragePrefix + "site"
	// CfgAccounted represents array of storages to be accounted
	CfgAccounted = cfgStoragePrefix + "accounted"
)
