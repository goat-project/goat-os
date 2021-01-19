package resource

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
	"github.com/gophercloud/gophercloud/pagination"
)

// FlavorReader structure for a Reader which read an array of flavors.
type FlavorReader struct {
}

// ReadResources reads an array of users.
func (ur *FlavorReader) ReadResources(client *gophercloud.ServiceClient) pagination.Pager {
	return flavors.ListDetail(client, flavors.ListOpts{})
}
