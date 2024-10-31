package gen

import (
	"context"
	"errors"
	"slices"

	"github.com/iancoleman/strcase"
	"github.com/nuzur/extension-sdk/client"
	"github.com/nuzur/extension-sdk/domainhelpers"
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
		return nil, errors.New("invalid request")
	}

	projectVersion := req.Deps.ProjectVersion

	entities := []SchemaEntity{}
	for i, e := range projectVersion.Entities {
		if slices.Contains(req.Configvalues.Entities, e.Uuid) {
			fields, indexes := MapEntityToTypes(e, req.Configvalues.DBType)
			selects := ResolveSelectStatements(e, req.Configvalues.DBType)
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
				SelectStatements: selects,
			}
			entities[i] = entityTemplate
		}
	}
	/*tpl := SchemaTemplate{
		Entities: entities,
	}*/

	return &GenerateResponse{
		DisplayBlocks:   []*pb.ExecutionResponseDisplayBlock{},
		FileDownloadUrl: "",
	}, nil
}
