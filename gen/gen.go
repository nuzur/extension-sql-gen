package gen

import (
	"context"

	"github.com/nuzur/extension-sdk/client"
	pb "github.com/nuzur/extension-sdk/idl/gen"
	"github.com/nuzur/extension-sql-gen/config"
)

type GenerateRequest struct {
	Configvalues *config.Values
	Deps         *client.BaseDependenciesResponse
}

type GenerateResponse struct {
	FileDownloadUrl string
	DisplayBlocks   []*pb.ExecutionResponseDisplayBlock
}

func Generate(ctx context.Context, req GenerateRequest) (*GenerateResponse, error) {
	if len(req.Configvalues.Entities) == 0 || len(req.Configvalues.Actions) == 0 {
		return nil, nil
	}
	return &GenerateResponse{}, nil
}
