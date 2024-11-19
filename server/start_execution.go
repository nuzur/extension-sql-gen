package server

import (
	"context"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/nuzur/extension-sdk/client"
	pb "github.com/nuzur/extension-sdk/idl/gen"
	nemgen "github.com/nuzur/extension-sdk/proto_deps/nem/idl/gen"
	"github.com/nuzur/extension-sql-gen/config"
	"github.com/nuzur/extension-sql-gen/gen"
	"golang.org/x/sync/errgroup"
)

func (s *server) StartExecution(ctx context.Context, req *pb.StartExecutionRequest) (*pb.StartExecutionResponse, error) {
	projectUUID := uuid.FromStringOrNil(req.ProjectUuid)
	projectVersionUUID := uuid.FromStringOrNil(req.ProjectVersionUuid)
	projectExtensionUUID := uuid.FromStringOrNil(req.ProjectExtensionUuid)

	start := time.Now()
	fmt.Printf("start exec! \n")
	configvalues := &config.Values{}
	err := s.client.ResolveConfigValues(ctx, client.ResolveConfigValuesRequest{
		ProjectUUID:          projectUUID,
		ProjectExtensionUUID: projectExtensionUUID,
		RawConfigValues:      req.ConfigValues,
	}, configvalues)
	if err != nil {
		return nil, err
	}

	fmt.Printf("config values: %v \n", time.Since(start))
	start = time.Now()

	eg, _ := errgroup.WithContext(ctx)

	var deps *client.BaseDependenciesResponse
	eg.Go(func() error {
		depsStart := time.Now()
		deps, err = s.client.GetBaseDependencies(ctx, client.BaseDependenciesRequest{
			ProjectUUID:        projectUUID,
			ProjectVersionUUID: projectVersionUUID,
		})
		if err != nil {
			return err
		}
		fmt.Printf("deps: %v \n", time.Since(depsStart))
		return nil
	})

	var exec *nemgen.ExtensionExecution
	eg.Go(func() error {
		metadata := gen.Metadata{
			ConfigValues: configvalues,
		}
		exec, err = s.client.CreateExecution(ctx, client.CreateExecutionRequest{
			ProjectUUID:          projectUUID,
			ProjectVersionUUID:   projectVersionUUID,
			ProjectExtensionUUID: projectExtensionUUID,
			Metadata:             metadata.ToString(),
		})
		if err != nil {
			return err
		}
		return nil
	})

	err = eg.Wait()
	if err != nil {
		return nil, err
	}

	fmt.Printf("deps + create exec: %v \n", time.Since(start))
	start = time.Now()

	res, err := gen.Generate(ctx, gen.GenerateRequest{
		ExecutionUUID: exec.Uuid,
		Client:        s.client,
		Configvalues:  configvalues,
		Deps:          deps,
	})
	if err != nil {
		s.client.UpdateExecution(ctx, client.UpdateExecutionRequest{
			ExecutionUUID:      uuid.FromStringOrNil(exec.Uuid),
			ProjectUUID:        projectUUID,
			ProjectVersionUUID: projectVersionUUID,
			Status:             pb.ExecutionStatus_EXECUTION_STATUS_FAILED,
			StatusMsg:          err.Error(),
		})
		return nil, err
	}
	fmt.Printf("gen: %v \n", time.Since(start))
	start = time.Now()

	// update final status
	_, err = s.client.UpdateExecution(ctx, client.UpdateExecutionRequest{
		ExecutionUUID:      uuid.FromStringOrNil(exec.Uuid),
		ProjectUUID:        projectUUID,
		ProjectVersionUUID: projectVersionUUID,
		Status:             pb.ExecutionStatus_EXECUTION_STATUS_SUCCEEDED,
		StatusMsg:          fmt.Sprintf("generated %d blocks and file: %s ", len(res.DisplayBlocks), res.FileDownloadUrl),
	})
	if err != nil {
		return nil, err
	}

	fmt.Printf("final update: %v \n", time.Since(start))

	return &pb.StartExecutionResponse{
		ExecutionUuid: exec.Uuid,
		Type:          pb.ExecutionResponseType_EXECUTION_RESPONSE_TYPE_FINAL,
		Data: &pb.ExecutionResponseTypeData{
			Final: &pb.ExecutionResponseTypeFinalData{
				Status:          pb.ExecutionStatus_EXECUTION_STATUS_SUCCEEDED,
				DisplayBlocks:   res.DisplayBlocks,
				FileDownloadUrl: res.FileDownloadUrl,
			},
		},
	}, nil
}
