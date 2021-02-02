package network

import (
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/floatingips"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
)

// NetUser represents "Resource" with information about project and his floating ips.
type NetUser struct {
	Project     *projects.Project
	FloatingIPs []floatingips.FloatingIP
}

// UnmarshalJSON function to implement Resource interface.
func (vnu *NetUser) UnmarshalJSON(b []byte) error {
	return vnu.Project.UnmarshalJSON(b)
}
