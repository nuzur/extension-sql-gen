package gen

import (
	nemgen "github.com/nuzur/extension-sdk/proto_deps/nem/idl/gen"
	"github.com/nuzur/extension-sql-gen/config"
)

func MapEntityToTypes(e *nemgen.Entity, dbType config.DBType) ([]SchemaField, []SchemaIndex) {
	fields := []SchemaField{}
	indexes := []SchemaIndex{}
	if e.Type != nemgen.EntityType_ENTITY_TYPE_STANDALONE {
		return fields, indexes
	}
	for _, f := range e.Fields {
		fieldType := ""
		if dbType == config.MYSQLDBType {
			fieldType = FieldTypeToMYSQL(f)
		}

		ft := SchemaField{
			Name: f.Identifier,
			Type: fieldType,
			Null: "NOT NULL",
		}

		if f.Unique {
			ft.Unique = "UNIQUE"
		}

		switch f.Type {
		case nemgen.FieldType_FIELD_TYPE_DATE:
			ft.Default = "default '2022-02-02'"
		case nemgen.FieldType_FIELD_TYPE_DATETIME:
			ft.Default = "default CURRENT_TIMESTAMP"
		}

		fields = append(fields, ft)

		for _, i := range e.TypeConfig.Standalone.Indexes {
			if len(i.Fields) == 1 && i.Fields[0].FieldUuid == f.Uuid && !f.Unique {
				indexes = append(indexes, SchemaIndex{
					Name:      f.Identifier,
					FieldName: f.Identifier,
				})
			}
		}
	}

	for i := 0; i < len(indexes)-1; i++ {
		indexes[i].HasComma = true
	}
	for i := 0; i < len(fields)-1; i++ {
		fields[i].HasComma = true
	}

	return fields, indexes
}
