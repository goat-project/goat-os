package server

import (
	"sync"

	"github.com/goat-project/goat-os/auth"
	"github.com/goat-project/goat-os/constants"
	"github.com/goat-project/goat-os/reader"
	"github.com/goat-project/goat-os/resource"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"

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

func (p *Processor) createReader(osClient *gophercloud.ProviderClient) {
	cClient, err := auth.CreateComputeV2ServiceClient(osClient)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("unable to create Compute V2 service client")
		return
	}

	p.reader = *reader.CreateReader(cClient)
}

// Reader gets reader.
func (p *Processor) Reader() *reader.Reader {
	return &p.reader
}

// Process provides listing of the servers with pagination.
func (p *Processor) Process(project projects.Project, osClient *gophercloud.ProviderClient, read chan resource.Resource,
	wg *sync.WaitGroup) {
	defer wg.Done()

	p.createReader(osClient)

	servs, err := p.reader.ListAllServers(project.ID)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error list servers")
		return
	}

	pages, err := servs.AllPages() // todo add openstack pagination and wg
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error get server pages")
		return
	}

	s, err := servers.ExtractServers(pages)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error extract servers")
		return
	}

	for i := range s {
		read <- &s[i]
	}
}

// RetrieveInfo calls method to retrieve server info.
func (p *Processor) RetrieveInfo(fullInfo chan resource.Resource, wg *sync.WaitGroup, vm resource.Resource) {
	defer wg.Done()

	if vm == nil {
		log.WithFields(log.Fields{}).Debug("retrieve info no vm")
		return
	}

	fullInfo <- vm
}
