package server

import (
	"context"

	pb "github.com/nuzur/extension-sdk/idl/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *server) StartExecution(context.Context, *pb.StartExecutionRequest) (*pb.StartExecutionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartExecution not implemented")
}
func (s *server) SubmitExectuionStep(context.Context, *pb.SubmitExectuionStepRequest) (*pb.SubmitExectuionStepResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SubmitExectuionStep not implemented")
}
func (s *server) GetExecutionStatus(context.Context, *pb.GetExecutionStatusRequest) (*pb.GetExecutionStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetExecutionStatus not implemented")
}
