package server

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/nuzur/extension-sdk/client"
	pb "github.com/nuzur/extension-sdk/idl/gen"
	sdkmapper "github.com/nuzur/extension-sdk/mapper"
	"github.com/nuzur/extension-sdk/proto_deps/nem/idl/gen"
	"github.com/nuzur/extension-sql-gen/constants"
	"github.com/nuzur/filetools"
)

func (s *server) GetExecution(ctx context.Context, req *pb.GetExecutionRequest) (*pb.GetExecutionResponse, error) {
	exec, err := s.client.GetExecution(ctx, uuid.FromStringOrNil(req.ExecutionUuid))
	if err != nil {
		return nil, err
	}

	// if the status is different that succeeded just return
	if exec.Status != gen.ExtensionExecutionStatus_EXTENSION_EXECUTION_STATUS_SUCCEEDED {
		return sdkmapper.MapExecutionToGetResponse(exec, nil, nil), nil
	}

	// if succeeded build the final result from the results

	// download results
	downloadRes, err := s.client.DownloadExecutionResults(ctx, client.DownloadExecutionResultsRequest{
		ExecutionUUID:      uuid.FromStringOrNil(req.ExecutionUuid),
		ProjectUUID:        uuid.FromStringOrNil(exec.ProjectUuid),
		ProjectVersionUUID: uuid.FromStringOrNil(exec.ProjectVersionUuid),
		FileExtension:      constants.ResultsFileExtension,
	})
	if err != nil || downloadRes == nil {
		return nil, err
	}

	// open the zip file
	read, err := zip.OpenReader(downloadRes.LocalFilePath)
	if err != nil {
		return nil, err
	}
	defer read.Close()

	finalData := &pb.ExecutionResponseTypeFinalData{
		Status:          pb.ExecutionStatus_EXECUTION_STATUS_SUCCEEDED,
		FileDownloadUrl: downloadRes.FileDownloadUrl,
		DisplayBlocks:   []*pb.ExecutionResponseDisplayBlock{},
	}
	// iterate through the files in the zip file
	for _, f := range read.File {
		// Open the current file
		v, err := f.Open()
		if err != nil {
			return nil, err
		}
		defer v.Close()

		// read file content
		var b bytes.Buffer
		_, err = io.Copy(&b, v)
		if err != nil {
			return nil, err
		}

		// build display blocks
		fileNameParts := strings.Split(f.Name, ".")
		identifier := strings.ReplaceAll(fileNameParts[0], "/", "")
		finalData.DisplayBlocks = append(finalData.DisplayBlocks, &pb.ExecutionResponseDisplayBlock{
			Identifier:  identifier,
			Title:       s.client.Localize(fmt.Sprintf("%s_title", identifier), req.Locale, identifier),
			Description: s.client.Localize(fmt.Sprintf("%s_description", identifier), req.Locale, identifier),
			Content:     b.String(),
			ContentType: pb.DisplayBlockContentType_DISPLAY_BLOCK_CONTENT_TYPE_SQL,
		})
	}

	// cleanup
	os.RemoveAll(path.Join(filetools.CurrentPath(), "previous-executions", fmt.Sprintf("%s.%s", req.ExecutionUuid, constants.ResultsFileExtension)))

	return sdkmapper.MapExecutionToGetResponse(exec, nil, finalData), nil
}
