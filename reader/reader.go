package reader

import (
	"time"

	"github.com/goat-project/goat-os/resource"
	networkReader "github.com/goat-project/goat-os/resource/network/reader"
	serverReader "github.com/goat-project/goat-os/resource/server/reader"
	storageReader "github.com/goat-project/goat-os/resource/storage/reader"

	log "github.com/sirupsen/logrus"

	"github.com/rafaeljesus/retry-go"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/users"
	"github.com/gophercloud/gophercloud/pagination"
)

// Reader structure to list resources and retrieve info for specific resource from Openstack.
type Reader struct {
	client *gophercloud.ServiceClient
}

type resourcesReaderI interface {
	ReadResources(*gophercloud.ServiceClient) pagination.Pager
}

type resourceReaderI interface {
	ReadResource(*gophercloud.ServiceClient) users.GetResult
}

const attempts = 3
const sleepTime = time.Second * 1

// CreateReader creates reader with gophercloud service client.
func CreateReader(client *gophercloud.ServiceClient) *Reader {
	if client == nil {
		log.WithFields(log.Fields{"error": "client is empty"}).Fatal("error create reader")
	}

	return &Reader{
		client: client,
	}
}

func (r *Reader) readResources(rri resourcesReaderI) (pagination.Pager, error) {
	var pager pagination.Pager
	var err error

	err = retry.Do(func() error {
		pager = rri.ReadResources(r.client)

		return err
	}, attempts, sleepTime)

	return pager, err
}

func (r *Reader) readResource(rri resourceReaderI) (users.GetResult, error) {
	var result users.GetResult
	var err error

	err = retry.Do(func() error {
		result = rri.ReadResource(r.client)

		return err
	}, attempts, sleepTime)

	return result, err
}

// ListAllServers lists all servers from Openstack.
func (r *Reader) ListAllServers(id string) (pagination.Pager, error) {
	return r.readResources(&serverReader.Servers{ProjectID: id})
}

// ListAllUsers lists all users from Openstack.
func (r *Reader) ListAllUsers() (pagination.Pager, error) {
	return r.readResources(&resource.UsersReader{})
}

// GetUser get user from Openstack.
func (r *Reader) GetUser(id string) (users.GetResult, error) {
	return r.readResource(&resource.UserReader{ID: id})
}

// ListAllFlavors lists all flavors from Openstack.
func (r *Reader) ListAllFlavors() (pagination.Pager, error) {
	return r.readResources(&resource.FlavorReader{})
}

// ListAllImages lists all images from Openstack.
func (r *Reader) ListAllImages(id string) (pagination.Pager, error) {
	return r.readResources(&storageReader.Image{ProjectID: id})
}

// ListAllShares lists all shares from Openstack.
func (r *Reader) ListAllShares(id string) (pagination.Pager, error) {
	return r.readResources(&storageReader.Share{ProjectID: id})
}

// ListAllVolumes lists all volumes.
func (r *Reader) ListAllVolumes(id string) (pagination.Pager, error) {
	return r.readResources(&storageReader.Volume{ProjectID: id})
}

// ListFloatingIPs lists all floating ips.
func (r *Reader) ListFloatingIPs() (pagination.Pager, error) {
	return r.readResources(&networkReader.FloatingIP{})
}

// ListAvailableProjects lists all available projects.
func (r *Reader) ListAvailableProjects() (pagination.Pager, error) {
	return r.readResources(&resource.ProjectReader{})
}
