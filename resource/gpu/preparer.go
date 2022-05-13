package gpu

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/goat-project/goat-os/constants"
	"github.com/goat-project/goat-os/initialize"
	"github.com/goat-project/goat-os/reader"
	"github.com/goat-project/goat-os/resource"
	"github.com/goat-project/goat-os/util"
	"github.com/goat-project/goat-os/writer"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"

	"golang.org/x/time/rate"
	"google.golang.org/grpc"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/spf13/viper"

	pb "github.com/goat-project/goat-proto-go"
	log "github.com/sirupsen/logrus"
)

// Preparer to prepare GPU data to specific structure for writing to Goat server.
type Preparer struct {
	identityReader reader.Reader
	computeReader  reader.Reader
	Writer         writer.Writer
	userIdentity   map[string]string
}

// CreatePreparer creates Preparer for virtual machine records.
func CreatePreparer(ir *reader.Reader, cr *reader.Reader, limiter *rate.Limiter, conn *grpc.ClientConn) *Preparer {
	if ir == nil {
		log.WithFields(log.Fields{}).Error(constants.ErrCreatePrepReaderNil)
		return nil
	}

	if cr == nil {
		log.WithFields(log.Fields{}).Error(constants.ErrCreatePrepReaderNil)
		return nil
	}

	if limiter == nil {
		log.WithFields(log.Fields{}).Error(constants.ErrCreatePrepLimiterNil)
		return nil
	}

	if conn == nil {
		log.WithFields(log.Fields{}).Error(constants.ErrCreatePrepConnNil)
		return nil
	}

	return &Preparer{
		identityReader: *ir,
		computeReader:  *cr,
		Writer:         *writer.CreateWriter(CreateWriter(limiter), conn),
	}
}

// InitializeMaps reads additional data for virtual machine record.
func (p *Preparer) InitializeMaps(wg *sync.WaitGroup) {
	defer wg.Done()

	wg.Add(1)
	go func() {
		defer wg.Done()
		p.userIdentity = initialize.UserIdentity(p.identityReader)
		if p.userIdentity == nil {
			log.WithFields(log.Fields{"error": "map is empty"}).Error("error create user identity map")
		}
	}()
}

// Preparation prepares GPU data for writing and call method to write.
func (p *Preparer) Preparation(acc resource.Resource, wg *sync.WaitGroup) {
	defer wg.Done()

	gpu := acc.(*Resource)
	if gpu == nil {
		log.WithFields(log.Fields{"error": "empty gpu"}).Error(constants.ErrPrepEmptyGPU)
		return
	}

	timeNow := time.Now()
	currentYear, currentMonth, _ := timeNow.Date()
	currentLocation := timeNow.Location()
	// unix of the first day of this month
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation).Unix()
	// accounting time - the first day of the month or the time when the server was created, in case it was created
	// during this month
	accountingTime := gpu.Server.Created.Unix()
	if firstOfMonth > accountingTime { // if the server was created this month
		accountingTime = firstOfMonth
	}
	// available duration is an active time during this month
	availableDuration := timeNow.Unix() - accountingTime
	if availableDuration < 0 { // should never happened
		availableDuration = 0
	}

	count, err := strconv.ParseFloat(gpu.ExtraSpecs["Accelerator:Number"], 32)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error convert gpu count")
	}

	scores, err := strconv.ParseFloat(gpu.ExtraSpecs["hw:cpu_cores"], 32)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error convert gpu cores")
	}

	cores := count * scores

	gpuRecord := pb.GPURecord{
		MeasurementMonth:     uint64(currentMonth),
		MeasurementYear:      uint64(currentYear),
		AssociatedRecordType: "cloud",
		AssociatedRecord:     gpu.Server.ID,
		GlobalUserName:       getGlobalUserName(p, gpu.Server),
		Fqan:                 gpu.Project.Name,
		SiteName:             getSiteName(),
		Count:                float32(count),
		Cores:                util.WrapUint32(fmt.Sprint(cores)),
		ActiveDuration:       util.WrapUint64(fmt.Sprint(availableDuration)),
		AvailableDuration:    uint64(availableDuration), // todo - uptime info from diagnostics v2.48
		//BenchmarkType: nil,
		//Benchmark: nil,
		Type:  gpu.ExtraSpecs["Accelerator:Type"],
		Model: util.WrapStr(gpu.ExtraSpecs["Accelerator:Model"]),
	}

	if err := p.Writer.Write(&gpuRecord); err != nil {
		log.WithFields(log.Fields{"error": err, "id": gpu.Server.ID}).Error(constants.ErrPrepWrite)
	}
}

// SendIdentifier sends identifier to Goat server.
func (p *Preparer) SendIdentifier() error {
	return p.Writer.SendIdentifier()
}

// Finish gets to know to the Goat server that a writing is finished and a response is expected.
// Then, it closes the gRPC connection.
func (p *Preparer) Finish() {
	p.Writer.Finish()

	log.WithFields(log.Fields{"type": "gpu"}).Debug("finished")
}

func getSiteName() string {
	siteName := viper.GetString(constants.CfgSiteName)
	if siteName == "" {
		log.WithFields(log.Fields{}).Error("no site name in configuration") // should never happen
	}

	return siteName
}

func getGlobalUserName(p *Preparer, server *servers.Server) *wrappers.StringValue {
	if p.userIdentity != nil {
		return util.WrapStr(p.userIdentity[server.UserID])
	}

	return nil
}
