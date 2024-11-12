package server

import (
	"context"

	pb "github.com/nuzur/extension-sdk/idl/gen"
)

func (s *server) GetMetadata(ctx context.Context, req *pb.GetMetadataRequest) (*pb.GetMetadataResponse, error) {
	initialMetadata := s.metadata

	switch req.Locale {
	case "es":
		initialMetadata.DisplayName = "Generador de SQL"
		initialMetadata.Description = "Genera el código SQL del proyecto."
	case "en":
		initialMetadata.DisplayName = "SQL Generator"
		initialMetadata.Description = "Generate the SQL code for the project."
	}

	return initialMetadata, nil
}
