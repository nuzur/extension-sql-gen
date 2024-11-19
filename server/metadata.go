package server

import (
	"context"

	pb "github.com/nuzur/extension-sdk/idl/gen"
)

func (s *server) GetMetadata(ctx context.Context, req *pb.GetMetadataRequest) (*pb.GetMetadataResponse, error) {
	initialMetadata := s.metadata

	initialMetadata.DisplayName = s.client.Localize("DisplayName", req.Locale, "SQL Generator")
	initialMetadata.Description = s.client.Localize("Description", req.Locale, "Generate SQL code for the project.")

	return initialMetadata, nil
}
