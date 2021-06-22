package reader

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumes"
	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
	"github.com/gophercloud/gophercloud/openstack/objectstorage/v1/containers"
	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/shares"
	"github.com/gophercloud/gophercloud/pagination"
)

// Image structure for a Reader which read an array of images.
type Image struct {
	ProjectID string
}

// Share structure for a Reader which read an array of shares.
type Share struct {
	ProjectID string
}

// Volume structure for a Reader which read an array of volumes.
type Volume struct {
	ProjectID string
}

// Swift structure for a Reader which read an array of swift containers.
type Swift struct {
}

// ReadResources reads an array of storages.
func (i *Image) ReadResources(client *gophercloud.ServiceClient) pagination.Pager {
	return images.List(client, images.ListOpts{Owner: i.ProjectID})
}

// ReadResources reads an array of storages.
func (s *Share) ReadResources(client *gophercloud.ServiceClient) pagination.Pager {
	return shares.ListDetail(client, shares.ListOpts{ProjectID: s.ProjectID})
}

// ReadResources reads an array of storages.
func (v *Volume) ReadResources(client *gophercloud.ServiceClient) pagination.Pager {
	return volumes.List(client, volumes.ListOpts{TenantID: v.ProjectID})
}

// ReadResources reads an array of storages.
func (s *Swift) ReadResources(client *gophercloud.ServiceClient) pagination.Pager {
	return containers.List(client, containers.ListOpts{Full: true})
}
