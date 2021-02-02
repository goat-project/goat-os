package network

import (
	"net"
	"sync"
	"time"

	"github.com/goat-project/goat-os/constants"
	"github.com/goat-project/goat-os/resource"
	"github.com/goat-project/goat-os/util"
	"github.com/goat-project/goat-os/writer"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/golang/protobuf/ptypes/wrappers"

	"github.com/spf13/viper"

	"golang.org/x/time/rate"
	"google.golang.org/grpc"

	pb "github.com/goat-project/goat-proto-go"
	log "github.com/sirupsen/logrus"
)

// Preparer to prepare network data to specific structure for writing to Goat server.
type Preparer struct {
	Writer writer.Writer
}

// CreatePreparer creates Preparer for network records.
func CreatePreparer(limiter *rate.Limiter, conn *grpc.ClientConn) *Preparer {
	if limiter == nil {
		log.WithFields(log.Fields{}).Error(constants.ErrCreatePrepLimiterNil)
		return nil
	}

	if conn == nil {
		log.WithFields(log.Fields{}).Error(constants.ErrCreatePrepConnNil)
		return nil
	}

	return &Preparer{
		Writer: *writer.CreateWriter(CreateWriter(limiter), conn),
	}
}

// InitializeMaps - only for VM relevant.
func (p *Preparer) InitializeMaps(wg *sync.WaitGroup) {
	defer wg.Done()
}

// Preparation prepares network data for writing and call method to write.
func (p *Preparer) Preparation(acc resource.Resource, wg *sync.WaitGroup) {
	defer wg.Done()

	netUser := acc.(*NetUser)
	if netUser.Project == nil {
		log.WithFields(log.Fields{}).Error(constants.ErrPrepEmptyNetUser)
		return
	}

	countIPv4, countIPv6 := countIPs(*netUser)

	if countIPv4 != 0 {
		ipv4Record := createIPRecord(*netUser, "IPv4", countIPv4)

		if err := p.Writer.Write(ipv4Record); err != nil {
			log.WithFields(log.Fields{"error": err, "id": netUser.Project.ID}).Error(constants.ErrPrepWrite)
		}
	}

	if countIPv6 != 0 {
		ipv6Record := createIPRecord(*netUser, "IPv6", countIPv6)

		if err := p.Writer.Write(ipv6Record); err != nil {
			log.WithFields(log.Fields{"error": err, "id": netUser.Project.ID}).Error(constants.ErrPrepWrite)
		}
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
	siteName := viper.GetString(constants.CfgNetworkSiteName)
	if siteName == "" {
		log.WithFields(log.Fields{}).Error(constants.ErrNoSiteName) // should never happen
	}

	return siteName
}

func getCloudComputeService() *wrappers.StringValue {
	return util.WrapStr(viper.GetString(constants.CfgNetworkCloudComputeService))
}

func getCloudType() string {
	ct := viper.GetString(constants.CfgNetworkCloudType)
	if ct == "" {
		log.WithFields(log.Fields{}).Error(constants.ErrNoCloudType) // should never happen
	}

	return ct
}

func countIPs(user NetUser) (uint32, uint32) {
	var countIPv4 uint32
	var countIPv6 uint32

	for _, fip := range user.FloatingIPs {
		ip := net.ParseIP(fip.IP)
		if ip == nil {
			continue
		}

		if ip4 := ip.To4(); ip4 != nil {
			countIPv4++
		} else if ip6 := ip.To16(); ip6 != nil {
			countIPv6++
		}
	}

	return countIPv4, countIPv6
}

func createIPRecord(netUser NetUser, ipType string, ipCount uint32) *pb.IpRecord {
	return &pb.IpRecord{
		MeasurementTime:     &timestamp.Timestamp{Seconds: time.Now().Unix()},
		SiteName:            getSiteName(),
		CloudComputeService: getCloudComputeService(),
		CloudType:           getCloudType(),
		LocalUser:           netUser.Project.ID,
		LocalGroup:          netUser.Project.DomainID,
		GlobalUserName:      netUser.Project.Name,
		Fqan:                "/" + netUser.Project.DomainID + "/Role=NULL/Capability=NULL",
		IpType:              ipType,
		IpCount:             ipCount,
	}
}
