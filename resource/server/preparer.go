package server

import (
	"net"
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

	"github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/spf13/viper"

	pb "github.com/goat-project/goat-proto-go"
	log "github.com/sirupsen/logrus"
)

// Preparer to prepare virtual machine data to specific structure for writing to Goat server.
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

// Preparation prepares virtual machine data for writing and call method to write.
func (p *Preparer) Preparation(acc resource.Resource, wg *sync.WaitGroup) {
	defer wg.Done()

	server := acc.(*SFStruct)
	if server == nil {
		log.WithFields(log.Fields{"error": "empty server"}).Error(constants.ErrPrepEmptyVM)
		return
	}

	sTime := util.WrapTime(&server.Server.Created)
	t := time.Now()
	eTime := util.WrapTime(&t) // todo get end time
	wallDuration := getWallDuration(sTime, eTime)

	var cpuCount uint32
	var memory *wrappers.UInt64Value
	var diskSize *wrappers.UInt64Value

	if server.Flavor != nil {
		cpuCount = uint32(server.Flavor.VCPUs)

		mem := uint64(server.Flavor.RAM)
		if mem != 0 {
			memory = &wrappers.UInt64Value{Value: mem}
		}

		disk := uint64(server.Flavor.Disk)
		if disk != 0 {
			diskSize = &wrappers.UInt64Value{Value: disk}
		}
	}

	serverRecord := pb.VmRecord{
		VmUuid:              server.Server.ID,
		SiteName:            getSiteName(),
		CloudComputeService: getCloudComputeService(),
		MachineName:         server.Server.Name,
		LocalUserId:         util.WrapStr(server.Server.UserID),
		LocalGroupId:        util.WrapStr(server.Server.TenantID),
		GlobalUserName:      getGlobalUserName(p, server.Server),
		Fqan:                getFqan(server.Server.TenantID),
		Status:              util.WrapStr(server.Server.Status),
		StartTime:           sTime,
		EndTime:             eTime,
		SuspendDuration:     getSuspendDuration(sTime, eTime, wallDuration),
		WallDuration:        wallDuration,
		CpuDuration:         getCPUDuration(wallDuration, cpuCount),
		CpuCount:            cpuCount,
		NetworkType:         nil,
		NetworkInbound:      nil, // todo?
		NetworkOutbound:     nil, // todo?
		PublicIpCount:       getPublicIPCount(server.Server),
		Memory:              memory,
		Disk:                diskSize,
		BenchmarkType:       nil, // todo?
		Benchmark:           nil, // todo?
		StorageRecordId:     nil,
		ImageId:             getImageID(server.Server),
		CloudType:           getCloudType(),
	}

	if err := p.Writer.Write(&serverRecord); err != nil {
		log.WithFields(log.Fields{"error": err, "id": server.Server.ID}).Error(constants.ErrPrepWrite)
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

	log.WithFields(log.Fields{"type": "server"}).Debug("finished")
}

func getSiteName() string {
	siteName := viper.GetString(constants.CfgSiteName)
	if siteName == "" {
		log.WithFields(log.Fields{}).Error("no site name in configuration") // should never happen
	}

	return siteName
}

func getCloudComputeService() *wrappers.StringValue {
	return util.WrapStr(viper.GetString(constants.CfgCloudComputeService))
}

func getGlobalUserName(p *Preparer, server *servers.Server) *wrappers.StringValue {
	if p.userIdentity != nil {
		return util.WrapStr(p.userIdentity[server.UserID])
	}

	return nil
}

func getFqan(tenantID string) *wrappers.StringValue {
	if tenantID != "" {
		return &wrappers.StringValue{Value: "/" + tenantID + "/Role=NULL/Capability=NULL"}
	}

	return nil
}

func getSuspendDuration(sTime, eTime *timestamp.Timestamp, wallDuration *duration.Duration) *duration.Duration {
	if eTime != nil && sTime != nil && wallDuration != nil {
		return &duration.Duration{Seconds: eTime.Seconds - sTime.Seconds - wallDuration.Seconds}
	}

	return nil
}

func getWallDuration(sTime, eTime *timestamp.Timestamp) *duration.Duration {
	if eTime != nil && sTime != nil {
		return &duration.Duration{Seconds: eTime.Seconds - sTime.Seconds}
	}

	return nil
} // todo should be /servers/{server_id}/diagnostics -> uptime

func getCPUDuration(wallDuration *duration.Duration, cpuCount uint32) *duration.Duration {
	if wallDuration != nil {
		return &duration.Duration{Seconds: wallDuration.Seconds * int64(cpuCount)}
	}

	return nil
}

func getPublicIPCount(server *servers.Server) *wrappers.UInt64Value {
	var sum int

	for _, a := range server.Addresses {
		for _, b := range a.([]interface{}) {
			address := b.(map[string]interface{})["addr"]
			if address != nil {
				ip := net.ParseIP(address.(string))
				if util.IsPublicIPv4(ip) {
					sum++
				}
			}
		}
	}

	if sum > 0 {
		return &wrappers.UInt64Value{Value: uint64(sum)}
	}

	return nil
}

func getImageID(server *servers.Server) *wrappers.StringValue {
	id := server.Image["id"]
	if id != nil {
		return util.WrapStr(id.(string))
	}

	return nil
}

func getCloudType() *wrappers.StringValue {
	ct := viper.GetString(constants.CfgCloudType)
	if ct == "" {
		log.WithFields(log.Fields{}).Error(constants.ErrNoCloudType) // should never happen
	}

	return &wrappers.StringValue{Value: ct}
}
