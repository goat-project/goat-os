package server

import (
	"sync"
	"time"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"

	"google.golang.org/grpc"

	"github.com/goat-project/goat-os/util"

	"github.com/goat-project/goat-os/resource"

	"github.com/goat-project/goat-os/writer"

	"golang.org/x/time/rate"

	"github.com/goat-project/goat-os/reader"

	"github.com/goat-project/goat-os/constants"

	"github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/spf13/viper"

	pb "github.com/goat-project/goat-proto-go"

	log "github.com/sirupsen/logrus"
)

// Preparer to prepare virtual machine data to specific structure for writing to Goat server.
type Preparer struct {
	reader reader.Reader
	Writer writer.Writer
}

// CreatePreparer creates Preparer for virtual machine records.
func CreatePreparer(reader *reader.Reader, limiter *rate.Limiter, conn *grpc.ClientConn) *Preparer {
	if reader == nil {
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
		reader: *reader,
		Writer: *writer.CreateWriter(CreateWriter(limiter), conn),
	}
}

// Preparation prepares virtual machine data for writing and call method to write.
func (p *Preparer) Preparation(acc resource.Resource, wg *sync.WaitGroup) {
	defer wg.Done()

	server := acc.(*servers.Server)
	// todo check that the server is correct server (not nil, has all needed attributes,...)
	id := server.ID

	sTime := util.WrapTime(&server.Created)
	t := time.Now()
	eTime := util.WrapTime(&t) // todo get end time
	wallDuration := getWallDuration(sTime, eTime)

	var cpuCount uint32
	vcpus := server.Flavor["vcpus"]
	if vcpus != nil {
		cpuCount = vcpus.(uint32)
	}

	serverRecord := pb.VmRecord{
		VmUuid:              server.ID,
		SiteName:            getSiteName(),
		CloudComputeService: getCloudComputeService(),
		MachineName:         server.Name,
		LocalUserId:         util.WrapStr(server.UserID),
		LocalGroupId:        util.WrapStr(server.TenantID),
		GlobalUserName:      util.WrapStr(""), // todo get from map of Users
		Fqan:                getFqan(server.TenantID),
		Status:              util.WrapStr(server.Status),
		StartTime:           sTime,
		EndTime:             eTime,
		SuspendDuration:     getSuspendDuration(sTime, eTime, wallDuration),
		WallDuration:        wallDuration,
		CpuDuration:         getCPUDuration(wallDuration, cpuCount),
		CpuCount:            cpuCount,
		NetworkType:         nil,
		NetworkInbound:      nil, // todo?
		NetworkOutbound:     nil, // todo?
		PublicIpCount:       getPublicIPCount(server),
		Memory:              getMemory(server),
		Disk:                getDiskSizes(server),
		BenchmarkType:       nil, // todo?
		Benchmark:           nil, // todo?
		StorageRecordId:     nil,
		ImageId:             getImageID(server),
		CloudType:           getCloudType(),
	}

	if err := p.Writer.Write(&serverRecord); err != nil {
		log.WithFields(log.Fields{"error": err, "id": id}).Error(constants.ErrPrepWrite)
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
	addresses := server.Addresses["public"]
	if addresses != nil {
		ads := addresses.([]servers.Address)
		length := len(ads)
		if length < 1 {
			return &wrappers.UInt64Value{Value: uint64(length)}
		}
	}

	return nil
}

func getMemory(server *servers.Server) *wrappers.UInt64Value {
	ram := server.Flavor["ram"]
	if ram != nil {
		return util.WrapUint64(ram.(string))
	}

	return nil
}

func getDiskSizes(server *servers.Server) *wrappers.UInt64Value {
	var disk *wrappers.UInt64Value
	var ephemeral *wrappers.UInt64Value

	d := server.Flavor["disk"]
	if d != nil {
		disk = util.WrapUint64(d.(string))
	}

	e := server.Flavor["ephemeral"]
	if e != nil {
		ephemeral = util.WrapUint64(e.(string))
	}

	if disk != nil && ephemeral != nil {
		return &wrappers.UInt64Value{Value: disk.Value + ephemeral.Value}
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
