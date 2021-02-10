package storage

import (
	"sync"
	"time"

	"github.com/goat-project/goat-os/initialize"

	"github.com/goat-project/goat-os/constants"
	"github.com/goat-project/goat-os/reader"
	"github.com/goat-project/goat-os/resource"
	"github.com/goat-project/goat-os/util"
	"github.com/goat-project/goat-os/writer"
	pb "github.com/goat-project/goat-proto-go"

	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/golang/protobuf/ptypes/wrappers"

	"golang.org/x/time/rate"
	"google.golang.org/grpc"

	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"

	"github.com/beevik/guid"
)

// Preparer to prepare storage data to specific structure for writing to Goat server.
type Preparer struct {
	reader       reader.Reader
	Writer       writer.Writer
	userIdentity map[string]string
}

// CreatePreparer creates Preparer for storage records.
func CreatePreparer(ir *reader.Reader, limiter *rate.Limiter, conn *grpc.ClientConn) *Preparer {
	if ir == nil {
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
		reader: *ir,
		Writer: *writer.CreateWriter(CreateWriter(limiter), conn),
	}
}

// InitializeMaps reads additional data for storage record.
func (p *Preparer) InitializeMaps(wg *sync.WaitGroup) {
	defer wg.Done()

	wg.Add(1)
	go func() {
		defer wg.Done()
		p.userIdentity = initialize.UserIdentity(p.reader)
	}()
}

// Preparation prepares storage data for writing and call method to write.
func (p *Preparer) Preparation(acc resource.Resource, wg *sync.WaitGroup) {
	defer wg.Done()

	storage := acc.(*images.Image)
	if storage == nil {
		log.WithFields(log.Fields{"error": "no image"}).Error(constants.ErrPrepEmptyImage)
		return
	}

	startTime := util.WrapTime(&storage.CreatedAt)
	now := time.Now().Unix()
	size := uint64(storage.SizeBytes)

	storageRecord := pb.StorageRecord{
		RecordID:      guid.New().String(),
		CreateTime:    &timestamp.Timestamp{Seconds: now},
		StorageSystem: viper.GetString(constants.CfgOpenstackIdentityEndpoint),
		Site:          util.WrapStr(viper.GetString(constants.CfgSite)),
		StorageShare:  util.WrapStr("image"), // todo rozlisit o aky typ ide - image, share (volume, swift)
		StorageMedia:  &wrappers.StringValue{Value: "disk"},
		// StorageClass: nil,
		FileCount: util.WrapStr(storage.File), // pre ostatne "1"
		// DirectoryPath: nil,
		LocalUser:    util.WrapStr(storage.Owner), // todo owner user id
		LocalGroup:   util.WrapStr(storage.Owner),
		UserIdentity: util.WrapStr(storage.Owner), // todo - owner's name
		Group:        util.WrapStr("/" + storage.Owner + "/Role=NULL/Capability=NULL"),
		// GroupAttribute: nil,
		// GroupAttributeType: nil,
		StartTime:                 startTime,
		EndTime:                   &timestamp.Timestamp{Seconds: now},
		ResourceCapacityUsed:      size,
		LogicalCapacityUsed:       &wrappers.UInt64Value{Value: size}, // todo - count
		ResourceCapacityAllocated: &wrappers.UInt64Value{Value: size}, // todo - count
	}

	if err := p.Writer.Write(&storageRecord); err != nil {
		log.WithFields(log.Fields{"error": err}).Error(constants.ErrPrepWrite)
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
