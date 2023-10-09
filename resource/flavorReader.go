// Package resource access
package resource

import (
	"github.com/goat-project/goat-os/result"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
	"github.com/gophercloud/gophercloud/pagination"
)

// FlavorReader structure for a Reader which read an array of flavors.
type FlavorReader struct {
}

// FlavorExtraSpecsReader structure for a Reader which read an array of flavor extra specs.
type FlavorExtraSpecsReader struct {
	FlavorID string
}

// ReadResources reads an array of users.
func (ur *FlavorReader) ReadResources(client *gophercloud.ServiceClient) pagination.Pager {
	return flavors.ListDetail(client, flavors.ListOpts{})
}

// ReadResource reads an array of flavor extra specs.
func (fesr *FlavorExtraSpecsReader) ReadResource(client *gophercloud.ServiceClient) result.Result {
	return flavors.ListExtraSpecs(client, fesr.FlavorID)
}
