package resource

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/users"
	"github.com/gophercloud/gophercloud/pagination"
)

// UserReader structure for a Reader which read an array of users.
type UserReader struct {
}

// ReadResources reads an array of users.
func (ur *UserReader) ReadResources(client *gophercloud.ServiceClient) pagination.Pager {
	return users.List(client, users.ListOpts{})
}
