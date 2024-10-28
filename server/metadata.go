package server

import (
	"context"

	pb "github.com/nuzur/extension-sdk/idl/gen"
)

func (s *server) GetMetadata(ctx context.Context, req *pb.GetMetadataRequest) (*pb.GetMetadataResponse, error) {
	return s.client.GetMetadata(ctx, req)
}
