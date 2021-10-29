package storage

import (
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumes"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
	"github.com/gophercloud/gophercloud/openstack/objectstorage/v1/containers"
	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/shares"
)

// SwiftContainer represents "Resource" with information about project and his container.
type SwiftContainer struct {
	Project   *projects.Project
	Container *containers.Container
}

// UnmarshalJSON function to implement Resource interface.
func (sc *SwiftContainer) UnmarshalJSON(b []byte) error {
	return sc.Project.UnmarshalJSON(b)
}

// PVolume represents "Resource" with information about project and his volume.
type PVolume struct {
	Project *projects.Project
	Volume  *volumes.Volume
}

// UnmarshalJSON function to implement Resource interface.
func (pv *PVolume) UnmarshalJSON(b []byte) error {
	return pv.Project.UnmarshalJSON(b)
}

// PShare represents "Resource" with information about project and his share.
type PShare struct {
	Project *projects.Project
	Share   *shares.Share
}

// UnmarshalJSON function to implement Resource interface.
func (ps *PShare) UnmarshalJSON(b []byte) error {
	return ps.Project.UnmarshalJSON(b)
}

// PImage represents "Resource" with information about project and his image.
type PImage struct {
	Project *projects.Project
	Image   *images.Image
}

// UnmarshalJSON function to implement Resource interface.
func (pi *PImage) UnmarshalJSON(b []byte) error {
	return pi.Project.UnmarshalJSON(b)
}
