package resource

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
	"github.com/gophercloud/gophercloud/pagination"
)

// ProjectReader structure for a Reader which reads an array of projects.
type ProjectReader struct {
}

// ReadResources reads an array of projects.
func (pr *ProjectReader) ReadResources(client *gophercloud.ServiceClient) pagination.Pager {
	return projects.ListAvailable(client)
}
