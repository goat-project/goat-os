package initialize

import (
	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/users"

	"github.com/goat-project/goat-os/reader"

	log "github.com/sirupsen/logrus"
)

// UserIdentity returns map of user ID and user name.
func UserIdentity(r reader.Reader) map[string]string {
	pages, err := r.ListAllUsers()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error list all users")
		return nil
	}

	u, err := pages.AllPages()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error list all user pages")
		return nil
	}

	usrs, err := users.ExtractUsers(u)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error extract all users")
		return nil
	}

	if len(usrs) == 0 {
		return nil
	}

	mUsers := make(map[string]string)

	for _, user := range usrs {
		if user.ID != "" {
			mUsers[user.ID] = user.Name
		}
	}

	return mUsers
}

// Flavor returns map of flavor ID and flavor structure.
func Flavor(r reader.Reader) map[string]*flavors.Flavor {
	pages, err := r.ListAllFlavors()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error list all flavors")
		return nil
	}

	f, err := pages.AllPages()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error list all flavor pages")
		return nil
	}

	flvrs, err := flavors.ExtractFlavors(f)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error extract all flavors")
		return nil
	}

	if len(flvrs) == 0 {
		return nil
	}

	mUsers := make(map[string]*flavors.Flavor)

	for i, flavor := range flvrs {
		if flavor.ID != "" {
			mUsers[flavor.ID] = &flvrs[i]
		}
	}

	return mUsers
}
