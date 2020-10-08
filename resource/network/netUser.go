package network

import (
	"github.com/gophercloud/gophercloud/openstack/identity/v3/users"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/floatingips"
)

// NetUser represents "Resource" with information about user and his floating ips.
type NetUser struct {
	User        *users.User
	FloatingIPs []floatingips.FloatingIP
}

// UnmarshalJSON function to implement Resource interface.
func (vnu *NetUser) UnmarshalJSON(b []byte) error {
	return vnu.User.UnmarshalJSON(b)
}
