// Package processor access
package processor

import (
	"github.com/goat-project/goat-os/resource"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
)

// Interface to process Resource data.
type Interface interface {
	ListProjects(chan projects.Project)
	ListResources(chan projects.Project, chan resource.Resource, gophercloud.AuthOptions)
}
