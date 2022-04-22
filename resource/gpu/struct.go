package gpu

import (
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
)

// GPUStruct represents "Resource" with information about project, server and his extra specs.
type GPUStruct struct {
	Project    *projects.Project
	Server     *servers.Server
	ExtraSpecs map[string]string
}

// UnmarshalJSON function to implement Resource interface.
func (gs *GPUStruct) UnmarshalJSON(b []byte) error {
	return gs.Server.UnmarshalJSON(b)
}
