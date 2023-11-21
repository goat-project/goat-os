// Package server access
package server

import (
	"sync"
	"time"

	"github.com/goat-project/goat-os/constants"
	"github.com/goat-project/goat-os/resource"

	"github.com/karrick/tparse/v2"

	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

// Filter contains times from/to filter records.
type Filter struct {
	recordsFrom time.Time
	recordsTo   time.Time
}

// CreateFilter creates Filter.
func CreateFilter() *Filter {
	recordsFrom := viper.GetTime(constants.CfgRecordsFrom)
	recordsTo := viper.GetTime(constants.CfgRecordsTo)

	periodStr := viper.GetString(constants.CfgRecordsForPeriod)
	period, err := tparse.AddDuration(time.Time{}, periodStr)
	if err != nil {
		log.WithFields(log.Fields{"period": periodStr}).Error("wrong format of period")
		period = time.Time{}
	}

	if (!recordsFrom.Equal(time.Time{}) || !recordsTo.Equal(time.Time{})) && !period.Equal(time.Time{}) {
		log.WithFields(log.Fields{
			"records-from": recordsFrom, "records-to": recordsTo, "period": periodStr,
		}).Fatal("cannot filter records from/to and records for a period in the same time")
	}

	if !period.Equal(time.Time{}) {
		now := time.Now()
		recFrom, err := tparse.AddDuration(now, "-"+periodStr)
		if err != nil {
			log.WithFields(log.Fields{"period": periodStr}).Error("wrong format of period")
		}

		log.WithFields(log.Fields{
			"record-from": recFrom, "record-to": now, "period": periodStr,
		}).Debug("filter set by a period")

		return &Filter{
			recordsFrom: recFrom,
			recordsTo:   now,
		}
	}

	if recordsTo.Equal(time.Time{}) {
		now := time.Now()

		log.WithFields(log.Fields{"record-from": recordsFrom, "record-to": now}).Debug("filter from a given time to now")

		return &Filter{
			recordsFrom: recordsFrom,
			recordsTo:   now,
		}
	}

	log.WithFields(log.Fields{"record-from": recordsFrom, "record-to": recordsTo}).Debug("filter set by times from and to")

	return &Filter{
		recordsFrom: recordsFrom,
		recordsTo:   recordsTo,
	}
}

// Filtering provides filtering given resources according to configuration or command line flags
// and writing to filtered channel.
func (f *Filter) Filtering(res resource.Resource, filtered chan resource.Resource, wg *sync.WaitGroup) {
	defer wg.Done()

	if res == nil {
		log.WithFields(log.Fields{"err": "no server"}).Error("error filter empty VM")
		return
	}

	server := res.(*SFStruct)

	stime := server.Server.Created
	etime := f.recordsTo // TODO server misses end time !!!

	// TODO server status contains only ACTIVE or IN_PROGESS,
	//  function list does not return inactive (error, deleted,
	//  etc.) servers. There should be used special call for
	//  deleted servers.
	//if server.Status != "ACTIVE" {
	//	etime = server.Updated
	//}

	if (stime.After(f.recordsFrom) || stime.Equal(f.recordsFrom)) &&
		(stime.Before(f.recordsTo) || stime.Equal(f.recordsTo)) &&
		(etime.After(f.recordsFrom) || etime.Equal(f.recordsFrom)) &&
		(etime.Before(f.recordsTo) || etime.Equal(f.recordsTo)) {
		filtered <- server
	}
}
