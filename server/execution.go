package server

import (
	"context"
	"encoding/json"

	"github.com/gofrs/uuid"
	"github.com/nuzur/extension-sdk/client"
	pb "github.com/nuzur/extension-sdk/idl/gen"
	"github.com/nuzur/extension-sql-gen/config"
	"github.com/nuzur/extension-sql-gen/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *server) StartExecution(ctx context.Context, req *pb.StartExecutionRequest) (*pb.StartExecutionResponse, error) {
	projectUUID := uuid.FromStringOrNil(req.ProjectUuid)
	projectVersionUUID := uuid.FromStringOrNil(req.ProjectVersionUuid)
	projectExtensionUUID := uuid.FromStringOrNil(req.ProjectExtensionUuid)

	configvalues, err := s.getConfigValues(ctx, client.ResolveConfigValuesRequest{
		ProjectUUID:          projectUUID,
		ProjectExtensionUUID: projectExtensionUUID,
		RawConfigValues:      req.ConfigValues,
	})
	if err != nil {
		return nil, err
	}

	deps, err := s.client.GetBaseDependencies(ctx, client.BaseDependenciesRequest{
		ProjectUUID:        projectUUID,
		ProjectVersionUUID: projectVersionUUID,
	})
	if err != nil {
		return nil, err
	}

	exec, err := s.client.CreateExecution(ctx, client.CreateExecutionRequest{
		ProjectUUID:          projectUUID,
		ProjectVersionUUID:   projectVersionUUID,
		ProjectExtensionUUID: projectExtensionUUID,
		Metadata:             "{}",
	})
	if err != nil {
		return nil, err
	}

	res, err := gen.Generate(ctx, gen.GenerateRequest{
		ExecutionUUID: exec.Uuid,
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

func (s *server) SubmitExectuionStep(context.Context, *pb.SubmitExectuionStepRequest) (*pb.SubmitExectuionStepResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SubmitExectuionStep not implemented")
}
func (s *server) GetExecution(context.Context, *pb.GetExecutionRequest) (*pb.GetExecutionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetExecutionStatus not implemented")
}

func (s *server) getConfigValues(ctx context.Context, req client.ResolveConfigValuesRequest) (*config.Values, error) {
	configValues, err := s.client.ResolveConfigValues(ctx, req)
	if err != nil {
		return nil, err
	}

	values := config.Values{}
	err = json.Unmarshal([]byte(*configValues), &values)
	if err != nil {
		return nil, err
	}

	return &values, nil
}
