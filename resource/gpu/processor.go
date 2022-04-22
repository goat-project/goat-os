package gpu

import (
	"fmt"
	"strings"
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

// Processor to process GPU's data.
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

// Process provides listing of the flavors, filtering flavors without `nvidia` in the name, listing of servers
// without pagination (all pages extracted in one step), listing of the extra specs for servers with flavor ID null or
// flavor name contained `nvidia`.
func (p *Processor) Process(project projects.Project, osClient *gophercloud.ProviderClient, read chan resource.Resource,
	wg *sync.WaitGroup) {
	defer wg.Done()

	p.createReader(osClient)

	flvrs, err := p.reader.ListAllFlavors()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error list flavor")
		return
	}

	pages, err := flvrs.AllPages() // todo add openstack pagination and wg
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error get flavor pages")
		return
	}

	allFlavors, err := flavors.ExtractFlavors(pages)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error extract servers")
		return
	}

	flavorsMap := make(map[string]*flavors.Flavor)
	containsNvidia := false

	if len(allFlavors) == 0 {
		log.WithFields(log.Fields{"error": err}).Error("error empty flavor list")
	} else {
		for i, flavor := range allFlavors {
			if flavor.ID != "" {
				if strings.Contains(allFlavors[i].Name, "nvidia") {
					containsNvidia = true
				}
				flavorsMap[flavor.ID] = &allFlavors[i]
			}
		}
	}

	if !containsNvidia {
		return // given project does not contain flavor with name `nvidia`, it does not support GPU
	}

	servs, err := p.reader.ListAllServers(project.ID)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error list servers")
		return
	}

	servsPages, err := servs.AllPages() // todo add openstack pagination and wg
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error get server pages")
		return
	}

	allServers, err := servers.ExtractServers(servsPages)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error extract servers")
		return
	}

	if len(allServers) < 1 {
		return // the project does not have any server
	}

	for i := range allServers {
		fid := allServers[i].Flavor["id"]

		var flavor *flavors.Flavor
		flavor = flavorsMap[fid.(string)]

		// condition takes only flavors with `nvidia` or broken flavors with old/null ID
		if flavor == nil || strings.Contains(flavor.Name, "nvidia") {
			eSpecs, err := p.reader.ListFlavorExtraSpecs(fmt.Sprint(fid))
			if err != nil {
				log.WithFields(log.Fields{"error": err}).Error("error list extra specs")
				continue
			}

			extraSpecs, err := eSpecs.(flavors.ListExtraSpecsResult).Extract()
			if err != nil {
				log.WithFields(log.Fields{"error": err}).Error("error extract servers")
				return
			}

			read <- &GPUStruct{Project: &project, Server: &allServers[i], ExtraSpecs: extraSpecs}
		}
	}
}
