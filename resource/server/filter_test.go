package server

import (
	"sync"
	"time"

	"github.com/goat-project/goat-os/resource"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"

	"github.com/goat-project/goat-os/constants"

	"github.com/spf13/viper"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Server Filter tests", func() {
	var (
		server *servers.Server
		wg     sync.WaitGroup
	)

	ginkgo.JustBeforeEach(func() {
		server = &servers.Server{Created: time.Unix(1540931164, 0)}

		viper.SetDefault(constants.CfgRecordsFrom, time.Time{})
		viper.SetDefault(constants.CfgRecordsTo, time.Time{})
		viper.SetDefault(constants.CfgRecordsForPeriod, time.Time{})
	})

	ginkgo.Describe("create filter", func() {
		ginkgo.Context("when no values are set", func() {
			ginkgo.It("should create filter with no restrictions", func() {
				filter := CreateFilter()

				gomega.Expect(filter.recordsFrom).To(gomega.Equal(time.Time{}))
				gomega.Expect(filter.recordsTo).To(gomega.And(
					gomega.BeTemporally("<", time.Now().Add(time.Minute)),
					gomega.BeTemporally(">", time.Now().Add(-time.Minute))))
			})
		})
	})

	ginkgo.Describe("create filter", func() {
		ginkgo.Context("when time from is set", func() {
			ginkgo.It("should create filter", func() {
				dateFrom := time.Now().Add(-48 * time.Hour)
				viper.SetDefault(constants.CfgRecordsFrom, dateFrom)

				filter := CreateFilter()

				gomega.Expect(filter.recordsFrom).To(gomega.Equal(dateFrom))
				gomega.Expect(filter.recordsTo).To(gomega.And(
					gomega.BeTemporally("<", time.Now().Add(time.Minute)),
					gomega.BeTemporally(">", time.Now().Add(-time.Minute))))
			})
		})
	})

	ginkgo.Describe("create filter", func() {
		ginkgo.Context("when time from and to are set", func() {
			ginkgo.It("should create filter", func() {
				dateFrom := time.Now().Add(-48 * time.Hour)
				dateTo := time.Now().Add(-24 * time.Hour)

				viper.SetDefault(constants.CfgRecordsFrom, dateFrom)
				viper.SetDefault(constants.CfgRecordsTo, dateTo)

				filter := CreateFilter()

				gomega.Expect(filter.recordsFrom).To(gomega.Equal(dateFrom))
				gomega.Expect(filter.recordsTo).To(gomega.Equal(dateTo))
			})
		})
	})

	ginkgo.Describe("create filter", func() {
		ginkgo.Context("when time to is set", func() {
			ginkgo.It("should create filter", func() {
				dateTo := time.Now().Add(-48 * time.Hour)
				viper.SetDefault(constants.CfgRecordsTo, dateTo)

				filter := CreateFilter()

				gomega.Expect(filter.recordsFrom).To(gomega.Equal(time.Time{}))
				gomega.Expect(filter.recordsTo).To(gomega.Equal(dateTo))
			})
		})
	})

	ginkgo.Describe("create filter", func() {
		ginkgo.Context("when time from and to and period are set", func() {
			ginkgo.It("should not create filter", func() {
				dateFrom := time.Now().Add(-48 * time.Hour)
				dateTo := time.Now().Add(-24 * time.Hour)
				period := "1y"

				viper.SetDefault(constants.CfgRecordsFrom, dateFrom)
				viper.SetDefault(constants.CfgRecordsTo, dateTo)
				viper.SetDefault(constants.CfgRecordsForPeriod, period)

				// TODO test Fatal error
			})
		})
	})

	ginkgo.Describe("create filter", func() {
		ginkgo.Context("when period is set", func() {
			ginkgo.It("should create filter", func() {
				period := "1y"
				viper.SetDefault(constants.CfgRecordsForPeriod, period)

				filter := CreateFilter()

				// handle leap year
				days := 365
				if isLeapYear(time.Now().Year()) || (isLeapYear(time.Now().Year()-1) && time.Now().Month() < 3) {
					days = 366
				}

				expectation := time.Now().Add(-time.Duration(days) * 24 * time.Hour)

				gomega.Expect(filter.recordsFrom).To(gomega.And(
					gomega.BeTemporally("<", expectation.Add(time.Minute)),
					gomega.BeTemporally(">", expectation.Add(-time.Minute))))

				gomega.Expect(filter.recordsTo).To(gomega.And(
					gomega.BeTemporally("<", time.Now().Add(time.Minute)),
					gomega.BeTemporally(">", time.Now().Add(-time.Minute))))
			})
		})
	})

	ginkgo.Describe("create filter", func() {
		ginkgo.Context("when period and time to are set", func() {
			ginkgo.It("should create filter", func() {
				dateTo := time.Now().Add(-24 * time.Hour)
				period := "1y"

				viper.SetDefault(constants.CfgRecordsForPeriod, period)
				viper.SetDefault(constants.CfgRecordsTo, dateTo)

				// TODO test Fatal error
			})
		})
	})

	ginkgo.Describe("create filter", func() {
		ginkgo.Context("when period and time from are set", func() {
			ginkgo.It("should create filter", func() {
				dateFrom := time.Now().Add(-48 * time.Hour)
				period := "1y"

				viper.SetDefault(constants.CfgRecordsForPeriod, period)
				viper.SetDefault(constants.CfgRecordsFrom, dateFrom)

				// TODO test Fatal error
			})
		})
	})

	ginkgo.Describe("filter virtual machine", func() {
		ginkgo.Context("when channel is empty and resource correct", func() {
			ginkgo.It("should not post vm to the channel", func(done ginkgo.Done) {
				// res := resources.CreateVirtualMachineWithID(1)
				filtered := make(chan resource.Resource)

				filter := CreateFilter()

				wg.Add(1)
				go filter.Filtering(server, filtered, &wg)

				gomega.Expect(filtered).To(gomega.BeEmpty())

				close(done)
			}, 0.2)
		})

		ginkgo.Context("when channel is empty and resource is not correct", func() {
			ginkgo.It("should not post vm to the channel", func(done ginkgo.Done) {
				filtered := make(chan resource.Resource)

				filter := CreateFilter()

				wg.Add(1)
				go filter.Filtering(nil, filtered, &wg)

				gomega.Expect(filtered).To(gomega.BeEmpty())

				close(done)
			}, 0.2)
		})

		// TODO add test with full channel

		ginkgo.Context("when channel is empty and resource time is in range", func() {
			ginkgo.It("should post vm to the channel", func(done ginkgo.Done) {
				dateTo := time.Now().Add(-24 * time.Hour)
				viper.SetDefault(constants.CfgRecordsTo, dateTo)

				filter := CreateFilter()

				filtered := make(chan resource.Resource)

				wg.Add(1)
				go filter.Filtering(server, filtered, &wg)

				gomega.Expect(<-filtered).To(gomega.Equal(server))

				close(done)
			}, 0.2)
		})

		ginkgo.Context("when channel is empty and resource time is out of range", func() {
			ginkgo.It("should not post vm to the channel", func(done ginkgo.Done) {
				dateTo := time.Now().Add(-2 * 356 * 24 * time.Hour)
				viper.SetDefault(constants.CfgRecordsTo, dateTo)

				filter := CreateFilter()

				filtered := make(chan resource.Resource)

				wg.Add(1)
				go filter.Filtering(server, filtered, &wg)

				gomega.Expect(filtered).To(gomega.BeEmpty())

				close(done)
			}, 0.2)
		})
	})
})

func isLeapYear(y int) bool {
	// convert int to Time - use the last day of the year, which is 31st December
	year := time.Date(y, time.December, 31, 0, 0, 0, 0, time.Local)
	days := year.YearDay()

	return days > 365
}
