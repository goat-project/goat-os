package resource

// Resource interface represents resource from Openstack.
type Resource interface {
	UnmarshalJSON([]byte) error // the only mutual function for server, user, image
}
