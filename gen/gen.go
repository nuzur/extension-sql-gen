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
	sdkgen "github.com/nuzur/extension-sdk/gen"
	pb "github.com/nuzur/extension-sdk/idl/gen"
	"github.com/nuzur/extension-sql-gen/config"
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

	genReq := &generateRequest{
		ExecutionUUID: req.ExecutionUUID,
		Configvalues:  configvalues,
		Data:          tpl,
		DisplayBlocks: &displayBlocks,
	}

	genFuncs := []func(context.Context, *generateRequest) error{
		generateCreates,
		generateInserts,
		generateUpdates,
		generateDeletes,
		generateSimpleSelects,
		generateIndexedSimpleSelects,
		generateIndexedCombinedSelects,
	}

	eg, _ := errgroup.WithContext(ctx)
	for _, genFunc := range genFuncs {
		eg.Go(func() error {
			return genFunc(ctx, genReq)
		})
	}
	err := eg.Wait()
	if err != nil {
		return nil, err
	}

	err = sdkgen.GenerateZip(ctx, sdkgen.ZipRequest{
		ExecutionUUID: req.ExecutionUUID,
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
		url, err := req.Client.UploadResults(ctx, client.UploadResultsRequest{
			ExecutionUUID:      uuid.FromStringOrNil(req.ExecutionUUID),
			ProjectUUID:        uuid.FromStringOrNil(req.Deps.Project.Uuid),
			ProjectVersionUUID: uuid.FromStringOrNil(req.Deps.ProjectVersion.Uuid),
			Data:               zipData,
			FileExtension:      "zip",
		})
		if err != nil || url == nil {
			return nil, err
		}
		downloadUrl = *url
	}

	// cleanup
	os.RemoveAll(path.Join("executions", req.ExecutionUUID))
	os.RemoveAll(path.Join("executions", fmt.Sprintf("%s.zip", req.ExecutionUUID)))

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
}

func generateCreates(ctx context.Context, req *generateRequest) error {
	if slices.Contains(req.Configvalues.Actions, config.CreateAction) {
		createdata, err := sdkgen.GenerateFile(ctx, sdkgen.FileRequest{
			ExecutionUUID:   req.ExecutionUUID,
			OutputFile:      "creates.sql",
			TemplateName:    fmt.Sprintf("%s_creates", req.Configvalues.DBType),
			Data:            req.Data,
			DisableGoFormat: true,
		})
		if err != nil {
			return err
		}
		req.mu.Lock()
		*req.DisplayBlocks = append(*req.DisplayBlocks, &pb.ExecutionResponseDisplayBlock{
			Identifier:  string(config.CreateAction),
			Title:       "Create SQL",
			Description: "",
			Content:     string(createdata),
		})
		req.mu.Unlock()
	}
	return nil
}

func generateInserts(ctx context.Context, req *generateRequest) error {
	if slices.Contains(req.Configvalues.Actions, config.InsertAction) {
		insertData, err := sdkgen.GenerateFile(ctx, sdkgen.FileRequest{
			ExecutionUUID:   req.ExecutionUUID,
			OutputFile:      "inserts.sql",
			TemplateName:    fmt.Sprintf("%s_inserts", req.Configvalues.DBType),
			Data:            req.Data,
			DisableGoFormat: true,
		})
		if err != nil {
			return err
		}
		req.mu.Lock()
		*req.DisplayBlocks = append(*req.DisplayBlocks, &pb.ExecutionResponseDisplayBlock{
			Identifier:  string(config.InsertAction),
			Title:       "Insert SQL",
			Description: "",
			Content:     string(insertData),
		})
		req.mu.Unlock()
	}
	return nil
}

func generateUpdates(ctx context.Context, req *generateRequest) error {
	if slices.Contains(req.Configvalues.Actions, config.UpdateAction) {
		updateData, err := sdkgen.GenerateFile(ctx, sdkgen.FileRequest{
			ExecutionUUID:   req.ExecutionUUID,
			OutputFile:      "updates.sql",
			TemplateName:    fmt.Sprintf("%s_updates", req.Configvalues.DBType),
			Data:            req.Data,
			DisableGoFormat: true,
		})
		if err != nil {
			return err
		}
		req.mu.Lock()
		*req.DisplayBlocks = append(*req.DisplayBlocks, &pb.ExecutionResponseDisplayBlock{
			Identifier:  string(config.UpdateAction),
			Title:       "Update SQL",
			Description: "",
			Content:     string(updateData),
		})
		req.mu.Unlock()
	}
	return nil
}

func generateDeletes(ctx context.Context, req *generateRequest) error {
	if slices.Contains(req.Configvalues.Actions, config.DeleteAction) {
		deleteData, err := sdkgen.GenerateFile(ctx, sdkgen.FileRequest{
			ExecutionUUID:   req.ExecutionUUID,
			OutputFile:      "deletes.sql",
			TemplateName:    fmt.Sprintf("%s_deletes", req.Configvalues.DBType),
			Data:            req.Data,
			DisableGoFormat: true,
		})
		if err != nil {
			return err
		}
		req.mu.Lock()
		*req.DisplayBlocks = append(*req.DisplayBlocks, &pb.ExecutionResponseDisplayBlock{
			Identifier:  string(config.DeleteAction),
			Title:       "Delete SQL",
			Description: "",
			Content:     string(deleteData),
		})
		req.mu.Unlock()
	}
	return nil
}

func generateSimpleSelects(ctx context.Context, req *generateRequest) error {
	if slices.Contains(req.Configvalues.Actions, config.SelectSimpleAction) {
		selectData, err := sdkgen.GenerateFile(ctx, sdkgen.FileRequest{
			ExecutionUUID:   req.ExecutionUUID,
			OutputFile:      "selects_simple.sql",
			TemplateName:    fmt.Sprintf("%s_selects_simple", req.Configvalues.DBType),
			Data:            req.Data,
			DisableGoFormat: true,
		})
		if err != nil {
			return err
		}
		req.mu.Lock()
		*req.DisplayBlocks = append(*req.DisplayBlocks, &pb.ExecutionResponseDisplayBlock{
			Identifier:  string(config.SelectSimpleAction),
			Title:       "Select Simple SQL",
			Description: "",
			Content:     string(selectData),
		})
		req.mu.Unlock()
	}
	return nil
}

func generateIndexedSimpleSelects(ctx context.Context, req *generateRequest) error {
	if slices.Contains(req.Configvalues.Actions, config.SelectForIndexedSimpleAction) {
		selectData, err := sdkgen.GenerateFile(ctx, sdkgen.FileRequest{
			ExecutionUUID:   req.ExecutionUUID,
			OutputFile:      "selects_indexed_simple.sql",
			TemplateName:    fmt.Sprintf("%s_selects_indexed_simple", req.Configvalues.DBType),
			Data:            req.Data,
			DisableGoFormat: true,
		})
		if err != nil {
			return err
		}
		req.mu.Lock()
		*req.DisplayBlocks = append(*req.DisplayBlocks, &pb.ExecutionResponseDisplayBlock{
			Identifier:  string(config.SelectForIndexedSimpleAction),
			Title:       "Select Indexed Simple SQL",
			Description: "",
			Content:     string(selectData),
		})
		req.mu.Unlock()
	}
	return nil
}

func generateIndexedCombinedSelects(ctx context.Context, req *generateRequest) error {
	if slices.Contains(req.Configvalues.Actions, config.SelectForIndexedCombinedAction) {
		selectData, err := sdkgen.GenerateFile(ctx, sdkgen.FileRequest{
			ExecutionUUID:   req.ExecutionUUID,
			OutputFile:      "selects_indexed_combined.sql",
			TemplateName:    fmt.Sprintf("%s_selects_indexed_combined", req.Configvalues.DBType),
			Data:            req.Data,
			DisableGoFormat: true,
		})
		if err != nil {
			return err
		}
		req.mu.Lock()
		*req.DisplayBlocks = append(*req.DisplayBlocks, &pb.ExecutionResponseDisplayBlock{
			Identifier:  string(config.SelectForIndexedCombinedAction),
			Title:       "Select Indexed Combined SQL",
			Description: "",
			Content:     string(selectData),
		})
		req.mu.Unlock()
	}
	return nil
}
