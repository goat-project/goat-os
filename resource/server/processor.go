package server

import (
	"sync"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"

	"github.com/goat-project/goat-os/constants"

	"github.com/goat-project/goat-os/reader"
	"github.com/goat-project/goat-os/resource"

	log "github.com/sirupsen/logrus"
)

// Processor to process server's data.
type Processor struct {
	reader reader.Reader
}

// CreateProcessor creates processor with reader.
func CreateProcessor(r *reader.Reader) *Processor {
	if r == nil {
		log.WithFields(log.Fields{}).Error(constants.ErrCreateProcReaderNil)
		return nil
	}

	return &Processor{
		reader: *r,
	}
}

// Process provides listing of the servers with pagination.
func (p *Processor) Process(read chan resource.Resource, wg *sync.WaitGroup) {
	defer wg.Done()

	servs, err := p.reader.ListAllServers()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Fatal("error list servers")
	}

	pages, err := servs.AllPages() // todo add openstack pagination and wg
	if err != nil {
		log.WithFields(log.Fields{"error": err /* todo add pageoffset */}).Fatal("error get server pages")
	}

	s, err := servers.ExtractServers(pages)
	if err != nil {
		log.WithFields(log.Fields{"error": err /* todo add pageoffset */}).Fatal("error extract servers")
	}

	for i := range s {
		read <- &s[i]
	}
}

// RetrieveInfo calls method to retrieve server info.
func (p *Processor) RetrieveInfo(fullInfo chan resource.Resource, wg *sync.WaitGroup, vm resource.Resource) {
	defer wg.Done()

	// todo do we need any other info for server?
	//id, err := vm.ID()
	//if err != nil {
	//	log.WithFields(log.Fields{"error": err}).Fatal("error get virtual machine id")
	//}
	//
	//v, err := p.reader.RetrieveVirtualMachineInfo(id)
	//if err != nil {
	//	log.WithFields(log.Fields{"error": err}).Fatal("error retrieve virtual machine info")
	//}
	//

	if vm == nil {
		log.WithFields(log.Fields{}).Debug("retrieve info no vm")
		return
	}

	fullInfo <- vm
}
