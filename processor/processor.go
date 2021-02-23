package processor

import (
	"sync"

	"github.com/goat-project/goat-os/reader"

	"github.com/goat-project/goat-os/auth"
	"github.com/gophercloud/gophercloud"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"

	"github.com/goat-project/goat-os/resource"
	log "github.com/sirupsen/logrus"
)

// Processor to process resource data.
type Processor struct {
	proc processorI
}

type processorI interface {
	Reader() *reader.Reader
	Process(projects.Project, *gophercloud.ProviderClient, chan resource.Resource, *sync.WaitGroup)
}

// CreateProcessor creates Processor to manage reading from OpenNebula.
func CreateProcessor(proc processorI) *Processor {
	return &Processor{
		proc: proc,
	}
}

// ListProjects lists projects from Openstack to create individual service clients.
func (p *Processor) ListProjects(projChan chan projects.Project) {
	availableProjects, err := p.proc.Reader().ListAvailableProjects()
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Fatal("unable to list available projects")
	}

	pages, err := availableProjects.AllPages()
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Fatal("unable to get pages of available projects")
	}

	projs, err := projects.ExtractProjects(pages)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Fatal("unable to extract available projects")
	}

	for i := range projs {
		projChan <- projs[i]
	}

	close(projChan)
}

// ListResources calls method to list resource from OpenNebula.
func (p *Processor) ListResources(projChan chan projects.Project, read chan resource.Resource,
	opts gophercloud.AuthOptions) {
	var wg sync.WaitGroup

	for project := range projChan {
		opts.TenantName = project.Name
		opts.Scope.ProjectName = project.Name

		osClient, err := auth.OpenstackClient(opts)
		if err != nil {
			log.WithFields(log.Fields{"err": err}).Error("unable to create Openstack client")
			return
		}

		wg.Add(1)
		go p.proc.Process(project, osClient, read, &wg)
	}

	wg.Wait()
	close(read)
}
