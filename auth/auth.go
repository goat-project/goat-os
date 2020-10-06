package auth

import (
	"github.com/goat-project/goat-os/constants"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/spf13/viper"
)

// OpenstackClient logs in to an OpenStack cloud found at the identity endpoint specified by the options,
// acquires a token, and returns a Provider Client instance that's ready to operate.
func OpenstackClient() (*gophercloud.ProviderClient, error) {
	return openstack.AuthenticatedClient(gophercloud.AuthOptions{
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
			DomainName:  viper.GetString(constants.CfgDomainName),
			System:      viper.GetBool(constants.CfgScopeSystem),
		},
		ApplicationCredentialID:     viper.GetString(constants.CfgAppCredentialID),
		ApplicationCredentialName:   viper.GetString(constants.CfgAppCredentialName),
		ApplicationCredentialSecret: viper.GetString(constants.CfgAppCredentialSecret),
	})
}

// CreateIdentityV3ServiceClient creates a ServiceClient that may be used to access the v3 identity service.
func CreateIdentityV3ServiceClient(client *gophercloud.ProviderClient) (*gophercloud.ServiceClient, error) {
	return openstack.NewIdentityV3(client, endpointOptions())
}

// CreateImageV2ServiceClient creates a ServiceClient that may be used to access the v2 image service.
func CreateImageV2ServiceClient(client *gophercloud.ProviderClient) (*gophercloud.ServiceClient, error) {
	return openstack.NewImageServiceV2(client, endpointOptions())
}

// CreateComputeV2ServiceClient creates a ServiceClient that may be used with the v2 compute package.
func CreateComputeV2ServiceClient(client *gophercloud.ProviderClient) (*gophercloud.ServiceClient, error) {
	return openstack.NewComputeV2(client, endpointOptions())
}

// CreateSharedFileSystemV2ServiceClient creates a ServiceClient that may be used with the v2 sharedFileSystem package.
func CreateSharedFileSystemV2ServiceClient(client *gophercloud.ProviderClient) (*gophercloud.ServiceClient, error) {
	return openstack.NewSharedFileSystemV2(client, endpointOptions())
}

func endpointOptions() gophercloud.EndpointOpts {
	return gophercloud.EndpointOpts{
		Type:         viper.GetString(constants.CfgEndpointType),
		Name:         viper.GetString(constants.CfgEndpointName),
		Region:       viper.GetString(constants.CfgEndpointRegion),
		Availability: availability(),
	}
}

func availability() gophercloud.Availability {
	switch viper.GetString(constants.CfgEndpointAvailability) {
	case "public":
		return gophercloud.AvailabilityPublic
	case "admin":
		return gophercloud.AvailabilityAdmin
	case "internal":
		return gophercloud.AvailabilityInternal
	default:
		return gophercloud.AvailabilityPublic
	}
}
