package storage

import (
	"sync"

	"github.com/goat-project/goat-os/auth"
	"github.com/goat-project/goat-os/constants"
	"github.com/goat-project/goat-os/reader"
	"github.com/goat-project/goat-os/resource"
	"github.com/goat-project/goat-os/util"
	"github.com/spf13/viper"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumes"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/shares"

	log "github.com/sirupsen/logrus"
)

const (
	image            = "image"
	sharedFileSystem = "sharedFileSystem"
	manila           = "manila"
	volume           = "volume"
	swiftContainer   = "swift"
	all              = "all"
)

// Processor to process storage data.
type Processor struct {
	computeReader      reader.Reader
	shareReader        reader.Reader
	blockStorageReader reader.Reader
}

// CreateProcessor creates Processor to manage reading from Openstack.
func CreateProcessor(r *reader.Reader) *Processor {
	if r == nil {
		log.WithFields(log.Fields{}).Error(constants.ErrCreateProcReaderNil)
		return nil
	}

	return &Processor{
		computeReader:      *r,
		shareReader:        *r,
		blockStorageReader: *r,
	}
}

func (p *Processor) createReader(osClient *gophercloud.ProviderClient, name string) {
	var client *gophercloud.ServiceClient
	var err error

	switch name {
	case image:
		client, err = auth.CreateComputeV2ServiceClient(osClient)
		if err != nil {
			log.WithFields(log.Fields{"err": err}).Error("unable to create Shared File System V2 service client")
			return
		}
		p.computeReader = *reader.CreateReader(client)
	case sharedFileSystem, manila:
		client, err = auth.CreateSharedFileSystemV2ServiceClient(osClient)
		if err != nil {
			log.WithFields(log.Fields{"err": err}).Error("unable to create Shared File System V2 service client")
			return
		}
		p.shareReader = *reader.CreateReader(client)
	case volume:
		client, err = auth.CreateNewBlockStorageV3ServiceClient(osClient)
		if err != nil {
			log.WithFields(log.Fields{"err": err}).Error("unable to create Shared File System V2 service client")
			return
		}
		p.blockStorageReader = *reader.CreateReader(client)
	}
}

// Reader gets reader.
func (p *Processor) Reader() *reader.Reader {
	return &p.computeReader
}

// Process provides listing of the images with pagination.
func (p *Processor) Process(project projects.Project, osClient *gophercloud.ProviderClient, read chan resource.Resource,
	wg *sync.WaitGroup) {
	defer wg.Done()

	id := project.ID

	accounted := viper.GetStringSlice(constants.CfgAccounted)

	if util.Contains(accounted, all) {
		wg.Add(3)
		go p.processImages(osClient, read, id, wg)
		go p.processShares(osClient, read, id, wg)
		go p.processVolumes(osClient, read, id, wg)
	} else {
		if util.Contains(accounted, image) {
			wg.Add(1)
			go p.processImages(osClient, read, id, wg)
		}

		if util.Contains(accounted, sharedFileSystem) || util.Contains(accounted, manila) {
			wg.Add(1)
			go p.processShares(osClient, read, id, wg)
		}

		if util.Contains(accounted, volume) {
			wg.Add(1)
			go p.processVolumes(osClient, read, id, wg)
		}
	}
}

func (p *Processor) processImages(osClient *gophercloud.ProviderClient, read chan resource.Resource, id string,
	wg *sync.WaitGroup) {
	defer wg.Done()

	p.createReader(osClient, image)

	imgs, err := p.computeReader.ListAllImages(id)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error list images")
		return
	}

	pages, err := imgs.AllPages() // todo add openstack pagination and wg
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error get image pages")
		return
	}

	s, err := images.ExtractImages(pages)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error extract images")
		return
	}

	for i := range s {
		read <- &s[i]
	}
}

func (p *Processor) processShares(osClient *gophercloud.ProviderClient, read chan resource.Resource, id string,
	wg *sync.WaitGroup) {
	defer wg.Done()

	p.createReader(osClient, sharedFileSystem)

	shrs, err := p.shareReader.ListAllShares(id)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error list shares")
		return
	}

	pages, err := shrs.AllPages() // todo add openstack pagination and wg
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error get share pages")
		return
	}

	s, err := shares.ExtractShares(pages)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error extract shares")
		return
	}

	for i := range s {
		read <- &s[i]
	}
}

func (p *Processor) processVolumes(osClient *gophercloud.ProviderClient, read chan resource.Resource, id string,
	wg *sync.WaitGroup) {
	defer wg.Done()

	p.createReader(osClient, volume)

	r, err := p.blockStorageReader.ListAllVolumes(id)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error list volumes")
		return
	}

	pages, err := r.AllPages() // todo add openstack pagination and wg
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error get volume pages")
		return
	}

	rs, err := volumes.ExtractVolumes(pages)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error extract volumes")
		return
	}

	for i := range rs {
		read <- &rs[i]
	}
}
