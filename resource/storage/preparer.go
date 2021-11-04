package storage

import (
	"strconv"
	"sync"
	"time"

	"github.com/goat-project/goat-os/constants"
	"github.com/goat-project/goat-os/initialize"
	"github.com/goat-project/goat-os/reader"
	"github.com/goat-project/goat-os/resource"
	"github.com/goat-project/goat-os/util"
	"github.com/goat-project/goat-os/writer"
	pb "github.com/goat-project/goat-proto-go"

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

	var storageRecord *pb.StorageRecord

	switch t := acc.(type) {
	case *PImage:
		storageRecord = prepareImage(t)
	case *PShare:
		storageRecord = prepareShare(t)
	case *PVolume:
		storageRecord = prepareVolume(t)
	case *SwiftContainer:
		storageRecord = prepareSwiftContainer(t)
	default:
		log.WithFields(log.Fields{"type": t}).Error("error unknown type")
	}

	if storageRecord == nil {
		log.WithFields(log.Fields{"error": "no storage record"}).Error(constants.ErrPrepEmptyImage)
		return
	}

	if err := p.Writer.Write(storageRecord); err != nil {
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

	log.WithFields(log.Fields{"type": "storage"}).Debug("finished")
}

func prepareImage(storage *PImage) *pb.StorageRecord {
	startTime := util.WrapTime(&storage.Image.CreatedAt)
	now := time.Now().Unix()
	size := uint64(storage.Image.SizeBytes)

	return &pb.StorageRecord{
		RecordID:      guid.New().String(),
		CreateTime:    &timestamp.Timestamp{Seconds: now},
		StorageSystem: viper.GetString(constants.CfgOpenstackIdentityEndpoint),
		Site:          util.WrapStr(viper.GetString(constants.CfgSite)),
		StorageShare:  util.WrapStr("image"),
		StorageMedia:  &wrappers.StringValue{Value: "disk"},
		// StorageClass: nil,
		FileCount: util.WrapStr(storage.Image.File),
		// DirectoryPath: nil,
		LocalUser:    util.WrapStr(storage.Image.Owner), // todo owner user id
		LocalGroup:   util.WrapStr(storage.Project.ID),
		UserIdentity: util.WrapStr(storage.Image.Owner), // todo - owner's name
		Group:        util.WrapStr(storage.Project.Name),
		// GroupAttribute: nil,
		// GroupAttributeType: nil,
		StartTime:                 startTime,
		EndTime:                   &timestamp.Timestamp{Seconds: now},
		ResourceCapacityUsed:      size,
		LogicalCapacityUsed:       &wrappers.UInt64Value{Value: size}, // todo - count
		ResourceCapacityAllocated: &wrappers.UInt64Value{Value: size}, // todo - count
	}
}

func prepareShare(storage *PShare) *pb.StorageRecord {
	startTime := util.WrapTime(&storage.Share.CreatedAt)
	now := time.Now().Unix()
	size := uint64(storage.Share.Size * 1024 * 1024 * 1024) // translate GB to bytes

	return &pb.StorageRecord{
		RecordID:      guid.New().String(),
		CreateTime:    &timestamp.Timestamp{Seconds: now},
		StorageSystem: viper.GetString(constants.CfgOpenstackIdentityEndpoint),
		Site:          util.WrapStr(viper.GetString(constants.CfgSite)),
		StorageShare:  util.WrapStr("share"),
		StorageMedia:  &wrappers.StringValue{Value: "disk"},
		// StorageClass: nil,
		FileCount: util.WrapStr("1"),
		// DirectoryPath: nil,
		LocalUser:    util.WrapStr(storage.Share.ProjectID), // todo owner user id
		LocalGroup:   util.WrapStr(storage.Project.ID),
		UserIdentity: util.WrapStr(storage.Share.ProjectID), // todo - owner's name
		Group:        util.WrapStr(storage.Project.Name),
		// GroupAttribute: nil,
		// GroupAttributeType: nil,
		StartTime:                 startTime,
		EndTime:                   &timestamp.Timestamp{Seconds: now},
		ResourceCapacityUsed:      size,
		LogicalCapacityUsed:       &wrappers.UInt64Value{Value: size}, // todo - count
		ResourceCapacityAllocated: &wrappers.UInt64Value{Value: size}, // todo - count
	}
}

func prepareVolume(storage *PVolume) *pb.StorageRecord {
	startTime := util.WrapTime(&storage.Volume.CreatedAt)
	now := time.Now().Unix()
	size := uint64(storage.Volume.Size * 1024 * 1024 * 1024) // translate GB to bytes

	return &pb.StorageRecord{
		RecordID:      guid.New().String(),
		CreateTime:    &timestamp.Timestamp{Seconds: now},
		StorageSystem: viper.GetString(constants.CfgOpenstackIdentityEndpoint),
		Site:          util.WrapStr(viper.GetString(constants.CfgSite)),
		StorageShare:  util.WrapStr("volume"),
		StorageMedia:  &wrappers.StringValue{Value: "disk"},
		// StorageClass: nil,
		FileCount: util.WrapStr("1"),
		// DirectoryPath: nil,
		LocalUser:    util.WrapStr(storage.Volume.UserID), // todo owner user id
		LocalGroup:   util.WrapStr(storage.Project.ID),
		UserIdentity: util.WrapStr(storage.Volume.UserID), // todo - owner's name
		Group:        util.WrapStr(storage.Project.Name),
		// GroupAttribute: nil,
		// GroupAttributeType: nil,
		StartTime:                 startTime,
		EndTime:                   &timestamp.Timestamp{Seconds: now},
		ResourceCapacityUsed:      size,
		LogicalCapacityUsed:       &wrappers.UInt64Value{Value: size}, // todo - count
		ResourceCapacityAllocated: &wrappers.UInt64Value{Value: size}, // todo - count
	}
}

func prepareSwiftContainer(storage *SwiftContainer) *pb.StorageRecord {
	// startTime := util.WrapTime(nil) // todo
	now := time.Now().Unix()
	size := uint64(storage.Container.Bytes)

	return &pb.StorageRecord{
		RecordID:      guid.New().String(),
		CreateTime:    &timestamp.Timestamp{Seconds: now},
		StorageSystem: viper.GetString(constants.CfgOpenstackIdentityEndpoint),
		Site:          util.WrapStr(viper.GetString(constants.CfgSite)),
		StorageShare:  util.WrapStr("swift"),
		StorageMedia:  &wrappers.StringValue{Value: "disk"},
		// StorageClass: nil,
		FileCount: &wrappers.StringValue{Value: strconv.FormatInt(storage.Container.Count, 10)},
		// DirectoryPath: nil,
		LocalUser:    util.WrapStr(storage.Project.ID),
		LocalGroup:   util.WrapStr(storage.Project.DomainID),
		UserIdentity: util.WrapStr(storage.Project.Name),
		Group:        util.WrapStr(storage.Project.Name),
		// GroupAttribute: nil,
		// GroupAttributeType: nil,
		StartTime:                 &timestamp.Timestamp{Seconds: now}, // todo //startTime,
		EndTime:                   &timestamp.Timestamp{Seconds: now},
		ResourceCapacityUsed:      size,
		LogicalCapacityUsed:       &wrappers.UInt64Value{Value: size}, // todo - count
		ResourceCapacityAllocated: &wrappers.UInt64Value{Value: size}, // todo - count
	}
}
