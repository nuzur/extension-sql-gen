package server

import (
	"context"

	"github.com/gofrs/uuid"
	pb "github.com/nuzur/extension-sdk/idl/gen"
	sdkmapper "github.com/nuzur/extension-sdk/mapper"
)

func (s *server) GetExecution(ctx context.Context, req *pb.GetExecutionRequest) (*pb.GetExecutionResponse, error) {
	exec, err := s.client.GetExecution(ctx, uuid.FromStringOrNil(req.ExecutionUuid))
	if err != nil {
		return nil, err
	}

	// TODO build step or final data based on the status
	return sdkmapper.MapExecutionToGetResponse(exec, nil, nil), nil
}
