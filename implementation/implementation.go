package implementation

import (
	"context"

	"github.com/nuzur/extension-sdk/client"
	pb "github.com/nuzur/extension-sdk/idl/gen"
	"go.uber.org/fx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	pb.UnimplementedNuzurExtensionServer
	client *client.Client
}

type Params struct {
	fx.In
	Client *client.Client
}

func New(params Params) (pb.NuzurExtensionServer, error) {
	return &server{
		client: params.Client,
	}, nil
}

func (s *server) GetMetadata(ctx context.Context, req *pb.GetMetadataRequest) (*pb.GetMetadataResponse, error) {
	return s.client.GetMetadata(ctx, req)
}
func (s *server) StartExecution(context.Context, *pb.StartExecutionRequest) (*pb.StartExecutionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartExecution not implemented")
}
func (s *server) SubmitExectuionStep(context.Context, *pb.SubmitExectuionStepRequest) (*pb.SubmitExectuionStepResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SubmitExectuionStep not implemented")
}
func (s *server) GetExecutionStatus(context.Context, *pb.GetExecutionStatusRequest) (*pb.GetExecutionStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetExecutionStatus not implemented")
}
