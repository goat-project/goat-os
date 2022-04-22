package gpu

import (
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
)

// Resource represents "GPU Resource" with information about project, server and his extra specs.
type Resource struct {
	Project    *projects.Project
	Server     *servers.Server
	ExtraSpecs map[string]string
}

// UnmarshalJSON function to implement Resource interface.
func (gs *Resource) UnmarshalJSON(b []byte) error {
	return gs.Server.UnmarshalJSON(b)
}
