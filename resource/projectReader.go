// Package resource access
package resource

import (
	"github.com/goat-project/goat-os/result"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
	"github.com/gophercloud/gophercloud/pagination"
	log "github.com/sirupsen/logrus"
)

// ProjectReader structure for a Reader which reads an array of projects.
type ProjectReader struct {
	ID string
}

// ReadResources reads an array of projects.
func (pr *ProjectReader) ReadResources(client *gophercloud.ServiceClient) pagination.Pager {
	return projects.ListAvailable(client)
}

// ReadResource reads a project by ID.
func (pr *ProjectReader) ReadResource(client *gophercloud.ServiceClient) result.Result {

	r, err := projects.Get(client, pr.ID).Extract()
	if err != nil {
		log.WithFields(log.Fields{"error": "GET request for project details failed"}).Fatal("error project info fetching")
	}
	return r
}
