package server

import (
	"fmt"
	"sync"

	"github.com/goat-project/goat-os/auth"
	"github.com/goat-project/goat-os/constants"
	"github.com/goat-project/goat-os/reader"
	"github.com/goat-project/goat-os/resource"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
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

	if len(s) < 1 {
		return
	}

	flavorsMap, err := p.listAllFlavors()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error list flavors")
	}

	if flavorsMap == nil {
		for i := range s {
			read <- &SFStruct{Server: &s[i], Flavor: nil}
		}

		return
	}

	for i := range s {
		var flavor *flavors.Flavor

		fid := s[i].Flavor["id"]
		if fid != nil {
			flavor = flavorsMap[fid.(string)]
		}
		read <- &SFStruct{Server: &s[i], Flavor: flavor}
	}
}

func (p *Processor) listAllFlavors() (map[string]*flavors.Flavor, error) {
	pages, err := p.reader.ListAllFlavors()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error list all flavors")
		return nil, err
	}

	f, err := pages.AllPages()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error list all flavor pages")
		return nil, err
	}

	flvrs, err := flavors.ExtractFlavors(f)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error extract all flavors")
		return nil, err
	}

	if len(flvrs) == 0 {
		return nil, fmt.Errorf("error empty flavor list")
	}

	flavorsMap := make(map[string]*flavors.Flavor)

	for i, flavor := range flvrs {
		if flavor.ID != "" {
			flavorsMap[flavor.ID] = &flvrs[i]
		}
	}

	return flavorsMap, nil
}
