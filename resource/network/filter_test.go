package network_test

import (
	"strconv"
	"sync"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/floatingips"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/users"

	"github.com/goat-project/goat-os/resource/network"

	"github.com/onsi/gomega"

	"github.com/goat-project/goat-os/resource"
	"github.com/onsi/ginkgo"
)

var _ = ginkgo.Describe("Network Filter tests", func() {
	var (
		filter   *network.Filter
		net      resource.Resource
		filtered chan resource.Resource
		wg       sync.WaitGroup
	)

	ginkgo.JustBeforeEach(func() {
		filter = network.CreateFilter()
		wg.Add(1)
	})

	ginkgo.Describe("filter network", func() {
		ginkgo.Context("when channel is empty and resource correct", func() {
			ginkgo.BeforeEach(func() {
				net = createTestNetwork(1)
				filtered = make(chan resource.Resource)
			})

			ginkgo.It("should post network to the channel", func(done ginkgo.Done) {
				go filter.Filtering(net, filtered, &wg)

				gomega.Expect(<-filtered).To(gomega.Equal(net))

				close(done)
			}, 0.2)
		})

		ginkgo.Context("when channel is empty and resource is not correct", func() {
			ginkgo.BeforeEach(func() {
				filtered = make(chan resource.Resource)
			})

			ginkgo.It("should not post network to the channel", func(done ginkgo.Done) {
				go filter.Filtering(nil, filtered, &wg)

				gomega.Expect(filtered).To(gomega.BeEmpty())

				close(done)
			}, 0.2)
		})
	})
})

func createTestNetwork(userID int) *network.NetUser {
	return &network.NetUser{
		User:        &users.User{ID: strconv.Itoa(userID)},
		FloatingIPs: []floatingips.FloatingIP{{ID: "1"}, {ID: "2"}},
	}
}
