package resource

import (
	"github.com/goat-project/goat-os/constants"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/users"
	"github.com/gophercloud/gophercloud/pagination"
	"github.com/spf13/viper"
)

// UsersReader structure for a Reader which read an array of users.
type UsersReader struct {
}

// UserReader structure for a Reader which read a user by ID.
type UserReader struct {
	ID string
}

// ReadResources reads an array of users.
func (ur *UsersReader) ReadResources(client *gophercloud.ServiceClient) pagination.Pager {
	return users.List(client, users.ListOpts{DomainID: viper.GetString(constants.CfgDomainID)})
}

// ReadResource reads a user by ID.
func (ur *UserReader) ReadResource(client *gophercloud.ServiceClient) users.GetResult {
	return users.Get(client, ur.ID)
}
