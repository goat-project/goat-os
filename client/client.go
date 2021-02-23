package client

import (
	"sync"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"

	"github.com/goat-project/goat-os/resource"

	"github.com/goat-project/goat-os/filter"
	"github.com/goat-project/goat-os/preparer"
	"github.com/goat-project/goat-os/processor"
)

// Client runs application.
type Client struct {
}

// Run reads, filters and writes Accountable.
func (c *Client) Run(processor processor.Interface, filter filter.Interface, preparer preparer.Interface,
	opts gophercloud.AuthOptions) {
	var mapWg sync.WaitGroup
	mapWg.Add(1)

	go preparer.InitializeMaps(&mapWg)

	// initialize channels
	projs := make(chan projects.Project)
	read := make(chan resource.Resource)
	filtered := make(chan resource.Resource)

	// create done channel
	done := make(chan bool)
	defer close(done)

	go processor.ListProjects(projs)

	go processor.ListResources(projs, read, opts)
	go filter.Filter(read, filtered)
	go preparer.Prepare(filtered, done, &mapWg)

	<-done
}
