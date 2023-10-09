package processor

import (
	"sync"

	"github.com/goat-project/goat-os/constants"
	"github.com/goat-project/goat-os/reader"
	"github.com/spf13/viper"

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
		if commonTagExists(getTags(), projs[i].Tags) {
			projChan <- projs[i]
		}
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

// Util
func getTags() []string {
	specifiedTags := []string{viper.GetString(constants.CfgDefaultTag)} // Default state for tags is a default tag

	if len(viper.GetStringSlice(constants.CfgTags)) != 0 { // There are some tags specified
		specifiedTags = []string{}
		specifiedTags = append(specifiedTags, viper.GetStringSlice(constants.CfgTags)...)
	}

	if viper.GetBool(constants.CfgIgnoreTags) { // Ignoring tags
		specifiedTags = []string{}
	}
	return specifiedTags
}

func isEmpty(arr []string) bool {
	return len(arr) == 0
}

func commonTagExists(specifiedTags, serverTags []string) bool {

	// Case when tags are ignored
	if isEmpty(specifiedTags) {
		return true
	}

	// Case when tags are not ignored but server has no tags
	if isEmpty(serverTags) {
		return false
	}

	// Create a map to store unique strings from the first array.
	stringMap := make(map[string]bool)

	// Populate the map with strings from the first array.
	for _, str := range specifiedTags {
		stringMap[str] = true
	}

	// Check if any string from the second array is present in the map.
	for _, str := range serverTags {
		if stringMap[str] {
			return true
		}
	}

	// If no common strings were found, return false.
	return false
}
