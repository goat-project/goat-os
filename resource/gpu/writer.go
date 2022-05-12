package gpu

import (
	"context"

	"github.com/goat-project/goat-os/writer"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/goat-project/goat-os/constants"
	"github.com/spf13/viper"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"

	pb "github.com/goat-project/goat-proto-go"

	log "github.com/sirupsen/logrus"
)

// Writer structure to write gpu data to Goat server.
type Writer struct {
	Stream      pb.AccountingService_ProcessGPUsClient
	rateLimiter *rate.Limiter
}

// CreateWriter creates Writer for gpu data.
func CreateWriter(limiter *rate.Limiter) *Writer {
	return &Writer{
		rateLimiter: limiter,
	}
}

// SetUp creates gRPC client and sets up Stream to process gpu data to Writer.
func (w *Writer) SetUp(conn *grpc.ClientConn) {
	// create grpc client
	grpcClient := pb.NewAccountingServiceClient(conn)

	// create Stream to process GPUs
	stream, err := grpcClient.ProcessGPUs(context.Background())
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Fatal("error create gRPC client stream")
	}

	w.Stream = stream
}

// Write writes GPU record to Goat server.
func (w *Writer) Write(record writer.Record) error {
	rec := record.(*pb.GPURecord)

	gpuData := &pb.GPUData{
		Data: &pb.GPUData_Gpu{
			Gpu: rec,
		},
	}

	return w.Stream.Send(gpuData)
}

// SendIdentifier sends identifier to Goat server.
func (w *Writer) SendIdentifier() error {
	gpuDataIdentifier := pb.GPUData_Identifier{Identifier: viper.GetString(constants.CfgIdentifier)}
	data := &pb.GPUData{
		Data: &gpuDataIdentifier,
	}

	return w.Stream.Send(data)
}

// Close gets to know to the goat server that a writing is finished and a response is expected.
func (w *Writer) Close() (*empty.Empty, error) {
	return w.Stream.CloseAndRecv()
}
