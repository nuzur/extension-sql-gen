package main_test

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/nuzur/extension-sdk/client"
	nemgen "github.com/nuzur/extension-sdk/proto_deps/nem/idl/gen"
	"github.com/nuzur/extension-sql-gen/config"
	"github.com/nuzur/extension-sql-gen/gen"
	"github.com/stretchr/testify/assert"
)

func TestGenMysql(t *testing.T) {
	pvdata, err := os.ReadFile("./testdata/project_version.json")
	assert.NoError(t, err)
	projectVerion := &nemgen.ProjectVersion{}
	err = json.Unmarshal(pvdata, projectVerion)
	assert.NoError(t, err)
	req := gen.GenerateRequest{
		ExecutionUUID: uuid.Must(uuid.NewV4()).String(),
		DisableUpload: true,
		Configvalues: &config.Values{
			DBType: config.MYSQLDBType,
			Entities: []string{
				"b8629dd5-f6e5-483f-893a-842357e171fc", "6f9ca9c7-6af3-4301-82d2-739ec84eab83",
			},
			Actions: []config.Action{
				config.CreateAction,
				config.DeleteAction,
				config.InsertAction,
				config.DeleteAction,
				config.SelectSimpleAction,
				config.SelectForIndexedSimpleAction,
				config.SelectForIndexedCombinedAction,
			},
		},
		Deps: &client.BaseDependenciesResponse{
			ProjectVersion: projectVerion,
		},
	}
	res, err := gen.Generate(context.Background(), req)
	assert.NoError(t, err)
	assert.NotNil(t, res)

	insertsData, err := os.ReadFile("./testdata/mysql_inserts.sql")
	assert.NoError(t, err)
	updatesData, err := os.ReadFile("./testdata/mysql_updates.sql")
	assert.NoError(t, err)
	deletesData, err := os.ReadFile("./testdata/mysql_deletes.sql")
	assert.NoError(t, err)
	createsData, err := os.ReadFile("./testdata/mysql_creates.sql")
	assert.NoError(t, err)
	selectsSimpleData, err := os.ReadFile("./testdata/mysql_selects_simple.sql")
	assert.NoError(t, err)
	selectsIndexedSimpleData, err := os.ReadFile("./testdata/mysql_selects_indexed_simple.sql")
	assert.NoError(t, err)
	selectsIndexedCombinedData, err := os.ReadFile("./testdata/mysql_selects_indexed_combined.sql")
	assert.NoError(t, err)

	for _, db := range res.DisplayBlocks {
		switch db.Identifier {
		case "insert":
			assert.Equal(t, string(insertsData), db.Content)
		case "update":
			assert.Equal(t, string(updatesData), db.Content)
		case "delete":
			assert.Equal(t, string(deletesData), db.Content)
		case "create":
			assert.Equal(t, string(createsData), db.Content)
		case "select-simple":
			assert.Equal(t, string(selectsSimpleData), db.Content)
		case "select-indexed-simple":
			assert.Equal(t, string(selectsIndexedSimpleData), db.Content)
		case "select-indexed-combined":
			assert.Equal(t, string(selectsIndexedCombinedData), db.Content)
		}
	}

}
