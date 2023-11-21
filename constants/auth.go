// Package constants access
package constants

// prefix for authentication options
const cfgAuthOptions = "auth-options."

// prefix for scope
const cfgScope = "scope-"

const (
	// CfgUsername represents username to authenticate with password
	CfgUsername = cfgAuthOptions + "username"
	// CfgUserID represents user id to authenticate with password
	CfgUserID = cfgAuthOptions + "user-id"
	// CfgPassword represents password to authenticate
	CfgPassword = cfgAuthOptions + "password"

	// CfgPasscode represents passcode to authenticate with TOTP
	CfgPasscode = cfgAuthOptions + "passcode"

	// CfgDomainID represents domain id to authenticate
	CfgDomainID = cfgAuthOptions + "domain-id"
	// CfgDomainName represents domain name to authenticate
	CfgDomainName = cfgAuthOptions + "domain-name"

	// CfgTenantID represents tenant id to authenticate
	CfgTenantID = cfgAuthOptions + "tenant-id"
	// CfgTenantName represents tenant name to authenticate
	CfgTenantName = cfgAuthOptions + "tenant-name"

	// CfgAllowReauth true to allow cache credentials in memory
	CfgAllowReauth = cfgAuthOptions + "allow-reauth"

	// CfgTokenID represents token TenantID to authenticate
	CfgTokenID = cfgAuthOptions + "token-id"

	// CfgScopeProjectID represents scope project TenantID
	CfgScopeProjectID = cfgAuthOptions + cfgScope + "project-id"
	// CfgScopeProjectName represents scope project name
	CfgScopeProjectName = cfgAuthOptions + cfgScope + "project-name"
	// CfgScopeDomainID represents scope domain TenantID
	CfgScopeDomainID = cfgAuthOptions + cfgScope + "domain-id"
	// CfgScopeDomainName represents scope domain name
	CfgScopeDomainName = cfgAuthOptions + cfgScope + "domain-name"
	// CfgScopeSystem represents scope system
	CfgScopeSystem = cfgAuthOptions + cfgScope + "system"

	// CfgAppCredentialID represents id to authenticate using application credentials
	CfgAppCredentialID = cfgAuthOptions + "application-credential-id"
	// CfgAppCredentialName represents name to authenticate using application credentials
	CfgAppCredentialName = cfgAuthOptions + "application-credential-name"
	// CfgAppCredentialSecret represents secret to authenticate using application credentials
	CfgAppCredentialSecret = cfgAuthOptions + "application-credential-secret" // nolint: gosec
)
