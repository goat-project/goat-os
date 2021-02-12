package reader

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/floatingips"
	"github.com/gophercloud/gophercloud/pagination"
)

// FloatingIP structure for a Reader which read floating IPs by tenant id.
type FloatingIP struct {
}

// ReadResources reads a server info.
func (r *FloatingIP) ReadResources(client *gophercloud.ServiceClient) pagination.Pager {
	return floatingips.List(client)
}
