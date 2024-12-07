package gen

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"slices"
	"sync"

	"github.com/gofrs/uuid"
	"github.com/iancoleman/strcase"
	"github.com/nuzur/extension-sdk/client"
	"github.com/nuzur/extension-sdk/domainhelpers"
	pb "github.com/nuzur/extension-sdk/idl/gen"
	"github.com/nuzur/extension-sql-gen/config"
	"github.com/nuzur/extension-sql-gen/constants"
	"github.com/nuzur/filetools"
	"golang.org/x/sync/errgroup"
)

type GenerateRequest struct {
	ExecutionUUID string
	Configvalues  *config.Values
	Client        *client.Client
	Deps          *client.BaseDependenciesResponse
	DisableUpload bool
}

type GenerateResponse struct {
	FileDownloadUrl string
	DisplayBlocks   []*pb.ExecutionResponseDisplayBlock
}

func Generate(ctx context.Context, req GenerateRequest) (*GenerateResponse, error) {
	configvalues := req.Configvalues
	if len(configvalues.Entities) == 0 || len(configvalues.Actions) == 0 {
		return nil, errors.New("invalid request")
	}

	projectVersion := req.Deps.ProjectVersion

	entities := []SchemaEntity{}
	for _, e := range projectVersion.Entities {
		if slices.Contains(configvalues.Entities, e.Uuid) {
			fields, indexes, constraints := MapEntityToTypes(e, projectVersion, configvalues.DBType)
			selects := ResolveSelectStatements(e, configvalues.DBType)
			primaryKeys := domainhelpers.EntityPrimaryKeys(e)
			primaryKeysIdentifiers := []string{}
			for _, pk := range primaryKeys {
				primaryKeysIdentifiers = append(primaryKeysIdentifiers, pk.Identifier)
			}
			entityTemplate := SchemaEntity{
				DBType:           req.Configvalues.DBType,
				Name:             e.Identifier,
				NameTitle:        strcase.ToCamel(e.Identifier),
				PrimaryKeys:      primaryKeysIdentifiers,
				Fields:           fields,
				Indexes:          indexes,
				Constraints:      constraints,
				SelectStatements: selects,
			}
			entities = append(entities, entityTemplate)
		}
	}
	tpl := SchemaTemplate{
		Entities: entities,
	}
	displayBlocks := []*pb.ExecutionResponseDisplayBlock{}

	eg, _ := errgroup.WithContext(ctx)
	for _, action := range configvalues.Actions {
		eg.Go(func() error {
			return generate(ctx, &generateRequest{
				ExecutionUUID: req.ExecutionUUID,
				Configvalues:  configvalues,
				Data:          tpl,
				DisplayBlocks: &displayBlocks,
				Action:        action,
			})
		})
	}
	err := eg.Wait()
	if err != nil {
		return nil, err
	}

	err = filetools.GenerateZip(ctx, filetools.ZipRequest{
		OutputPath: "executions",
		Identifier: req.ExecutionUUID,
	})
	if err != nil {
		return nil, err
	}

	zipData, err := os.ReadFile(path.Join("executions", fmt.Sprintf("%s.zip", req.ExecutionUUID)))
	if err != nil {
		return nil, err
	}

	downloadUrl := ""
	if !req.DisableUpload {
		url, err := req.Client.UploadExecutionResults(ctx, client.UploadResultsRequest{
			ExecutionUUID:      uuid.FromStringOrNil(req.ExecutionUUID),
			ProjectUUID:        uuid.FromStringOrNil(req.Deps.Project.Uuid),
			ProjectVersionUUID: uuid.FromStringOrNil(req.Deps.ProjectVersion.Uuid),
			Data:               zipData,
			FileExtension:      constants.ResultsFileExtension,
		})

		if err != nil || url == nil {
			return nil, err
		}

		downloadUrl = *url
	}

	// cleanup
	os.RemoveAll(path.Join(filetools.CurrentPath(), "executions", req.ExecutionUUID))
	os.RemoveAll(path.Join(filetools.CurrentPath(), "executions", fmt.Sprintf("%s.zip", req.ExecutionUUID)))

	return &GenerateResponse{
		DisplayBlocks:   displayBlocks,
		FileDownloadUrl: downloadUrl,
	}, nil
}

type generateRequest struct {
	mu            sync.Mutex
	ExecutionUUID string
	Configvalues  *config.Values
	Data          SchemaTemplate
	DisplayBlocks *[]*pb.ExecutionResponseDisplayBlock
	Action        config.Action
}

func generate(ctx context.Context, req *generateRequest) error {
	data, err := filetools.GenerateFile(ctx, filetools.FileRequest{
		OutputPath:      path.Join("executions", req.ExecutionUUID, fmt.Sprintf("%s.sql", string(req.Action))),
		TemplateName:    fmt.Sprintf("%s_%s", string(req.Action), req.Configvalues.DBType),
		Data:            req.Data,
		DisableGoFormat: true,
	})
	if err != nil {
		return err
	}
	req.mu.Lock()
	*req.DisplayBlocks = append(*req.DisplayBlocks, &pb.ExecutionResponseDisplayBlock{
		Identifier:  string(req.Action),
		Content:     string(data),
		ContentType: pb.DisplayBlockContentType_DISPLAY_BLOCK_CONTENT_TYPE_SQL,
	})
	req.mu.Unlock()

	return nil
}
