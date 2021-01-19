package processor

import (
	"github.com/goat-project/goat-os/resource"
)

// Interface to process Resource data.
type Interface interface {
	ListResources(chan resource.Resource)
	RetrieveInfoResource(chan resource.Resource, chan resource.Resource)
}
