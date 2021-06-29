package server

import (
	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
)

// SFStruct represents "Resource" with information about server and his flavor.
type SFStruct struct {
	Server *servers.Server
	Flavor *flavors.Flavor
}

// UnmarshalJSON function to implement Resource interface.
func (sf *SFStruct) UnmarshalJSON(b []byte) error {
	return sf.Server.UnmarshalJSON(b)
}
