package server

import (
	"context"

	"github.com/gofrs/uuid"
	pb "github.com/nuzur/extension-sdk/idl/gen"
	sdkmapper "github.com/nuzur/extension-sdk/mapper"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *server) GetExecution(ctx context.Context, req *pb.GetExecutionRequest) (*pb.GetExecutionResponse, error) {
	exec, err := s.client.GetExecution(ctx, uuid.FromStringOrNil(req.ExecutionUuid))
	if err != nil {
		return nil, err
	}

	if exec.ExtensionUuid == s.metadata.Uuid {
		// TODO build step or final data based on the status
		return sdkmapper.MapExecutionToGetResponse(exec, nil, nil), nil
	}

	return nil, status.Errorf(codes.InvalidArgument, "execution not found")
}
