package network

import (
	"sync"

	"github.com/goat-project/goat-os/constants"

	"github.com/goat-project/goat-os/auth"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/floatingips"

	"github.com/goat-project/goat-os/reader"
	"github.com/goat-project/goat-os/resource"

	log "github.com/sirupsen/logrus"
)

// Processor to process network data.
type Processor struct {
	reader reader.Reader
}

// CreateProcessor creates Processor to manage reading from Openstack.
func CreateProcessor(r *reader.Reader) *Processor {
	if r == nil {
		log.WithFields(log.Fields{}).Error(constants.ErrCreateProcReaderNil)
		return nil
	}

	return &Processor{
		reader: *r,
	}
}

// Reader gets reader.
func (p *Processor) Reader() *reader.Reader {
	return &p.reader
}

func (p *Processor) createReader(osClient *gophercloud.ProviderClient) {
	cClient, err := auth.CreateComputeV2ServiceClient(osClient)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("unable to create Compute V2 service client")
		return
	}

	p.reader = *reader.CreateReader(cClient)
}

// Process provides listing of the users.
func (p *Processor) Process(project projects.Project, osClient *gophercloud.ProviderClient, read chan resource.Resource,
	wg *sync.WaitGroup) {
	defer wg.Done()

	p.createReader(osClient)

	fips, err := p.reader.ListFloatingIPs()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error retrieve floating Ips")
		return
	}

	pages, err := fips.AllPages() // todo add openstack pagination and wg
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error get fip pages")
		return
	}

	floatings, err := floatingips.ExtractFloatingIPs(pages)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error extract floating ips")
		return
	}

	read <- &NetUser{
		Project:     &project,
		FloatingIPs: floatings,
	}
}

// RetrieveInfo about virtual machines specific for a given user.
func (p *Processor) RetrieveInfo(fullInfo chan resource.Resource, wg *sync.WaitGroup, res resource.Resource) {
	defer wg.Done()

	fullInfo <- res
}
