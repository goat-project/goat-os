package gpu

import (
	"sync"

	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"

	"github.com/goat-project/goat-os/resource"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("GPU Filter tests", func() {
	var (
		filter   *Filter
		res      resource.Resource
		filtered chan resource.Resource
		wg       sync.WaitGroup
	)

	ginkgo.JustBeforeEach(func() {
		filter = CreateFilter()
		wg.Add(1)
	})

	ginkgo.Describe("filter gpu", func() {
		ginkgo.Context("when channel is empty and resource correct", func() {
			ginkgo.BeforeEach(func() {
				res = &images.Image{ID: "1"}
				filtered = make(chan resource.Resource)
			})

			ginkgo.It("should post gpu to the channel", func(done ginkgo.Done) {
				go filter.Filtering(res, filtered, &wg)

				gomega.Expect(<-filtered).To(gomega.Equal(res))

				close(done)
			}, 0.2)
		})

		ginkgo.Context("when channel is empty and resource is not correct", func() {
			ginkgo.BeforeEach(func() {
				filtered = make(chan resource.Resource)
			})

			ginkgo.It("should not post gpu to the channel", func(done ginkgo.Done) {
				go filter.Filtering(nil, filtered, &wg)

				gomega.Expect(filtered).To(gomega.BeEmpty())

				close(done)
			}, 0.2)
		})
	})
})
