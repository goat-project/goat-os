package storage

import (
	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
	"github.com/gophercloud/gophercloud/openstack/objectstorage/v1/containers"
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
