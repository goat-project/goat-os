package reader

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
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

// ReadResources reads an array of storages.
func (i *Image) ReadResources(client *gophercloud.ServiceClient) pagination.Pager {
	return images.List(client, images.ListOpts{Owner: i.ProjectID})
}

// ReadResources reads an array of storages.
func (s *Share) ReadResources(client *gophercloud.ServiceClient) pagination.Pager {
	return shares.ListDetail(client, shares.ListOpts{ProjectID: s.ProjectID})
}
