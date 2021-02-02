package reader

import (
	"github.com/gophercloud/gophercloud/pagination"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
)

// Servers structure for a Reader which reads an array of servers.
type Servers struct {
}

// ReadResources reads servers.
func (r *Servers) ReadResources(client *gophercloud.ServiceClient) pagination.Pager {
	return servers.List(client, nil)
}
