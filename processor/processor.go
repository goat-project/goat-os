package processor

import (
	"sync"

	"github.com/goat-project/goat-os/resource"
	log "github.com/sirupsen/logrus"
)

// Processor to process resource data.
type Processor struct {
	proc processorI
}

type processorI interface {
	Process(chan resource.Resource, *sync.WaitGroup)
	RetrieveInfo(chan resource.Resource, *sync.WaitGroup, resource.Resource)
}

// CreateProcessor creates Processor to manage reading from OpenNebula.
func CreateProcessor(proc processorI) *Processor {
	return &Processor{
		proc: proc,
	}
}

// ListResources calls method to list resource from OpenNebula.
func (p *Processor) ListResources(read chan resource.Resource) {
	var wg sync.WaitGroup

	wg.Add(1)
	go p.proc.Process(read, &wg)

	wg.Wait()
	close(read)
}

// RetrieveInfoResource range over filtered resource and calls method to retrieve resource info.
func (p *Processor) RetrieveInfoResource(filtered, fullInfo chan resource.Resource) {
	var wg sync.WaitGroup

	for accountable := range filtered {
		if accountable == nil {
			log.WithFields(log.Fields{"error": "no accountable"}).Error("error retrieve resource info")
			continue
		}

		wg.Add(1)
		go p.proc.RetrieveInfo(fullInfo, &wg, accountable)
	}

	wg.Wait()
	close(fullInfo)
}
