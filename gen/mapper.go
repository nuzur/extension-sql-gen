package gen

import (
	"sort"

	nemgen "github.com/nuzur/extension-sdk/proto_deps/nem/idl/gen"
	"github.com/nuzur/extension-sql-gen/config"
)

func MapEntityToTypes(e *nemgen.Entity, dbType config.DBType) ([]SchemaField, []SchemaIndex) {
	fields := []SchemaField{}
	indexes := []SchemaIndex{}
	if e.Type != nemgen.EntityType_ENTITY_TYPE_STANDALONE {
		return fields, indexes
	}

	identifers := make(map[string]string)
	for _, f := range e.Fields {
		if f.Status == nemgen.FieldStatus_FIELD_STATUS_ACTIVE {
			identifers[f.Uuid] = f.Identifier
			ft := mapField(f, dbType)
			fields = append(fields, ft)
		}
	}

	if e.TypeConfig != nil && e.TypeConfig.Standalone != nil {
		for _, i := range e.TypeConfig.Standalone.Indexes {
			if i.Status == nemgen.IndexStatus_INDEX_STATUS_ACTIVE {
				fieldNames := make(map[string]string)
				for _, fi := range i.Fields {
					fieldNames[fi.FieldUuid] = identifers[fi.FieldUuid]
				}

				indexTypePrefix := ""
				if i.Type == nemgen.IndexType_INDEX_TYPE_UNIQUE {
					indexTypePrefix = "UNIQUE "
				}
				if i.Type == nemgen.IndexType_INDEX_TYPE_FULLTEXT {
					indexTypePrefix = "FULLTEXT "
				}

				indexType := ""
				indexTypeSort := 0

				switch i.Type {
				case nemgen.IndexType_INDEX_TYPE_UNIQUE:
					indexType = "unique"
					indexTypeSort = 2
				case nemgen.IndexType_INDEX_TYPE_PRIMARY:
					indexType = "primary"
					indexTypeSort = 0
				case nemgen.IndexType_INDEX_TYPE_INDEX:
					indexType = "index"
					indexTypeSort = 1
				case nemgen.IndexType_INDEX_TYPE_FULLTEXT:
					indexType = "fulltext"
					indexTypeSort = 3
				}

				indexes = append(indexes, SchemaIndex{
					Name:       i.Identifier,
					Index:      i,
					FieldNames: fieldNames,
					Type:       indexType,
					TypeSort:   indexTypeSort,
					TypePrefix: indexTypePrefix,
				})

				sort.Slice(indexes, func(i, j int) bool {
					return indexes[i].TypeSort < indexes[j].TypeSort
				})
			}
		}
	}

	if len(indexes) > 0 {
		for i := 0; i < len(indexes)-1; i++ {
			indexes[i].HasComma = true
		}
	}

	if len(fields) > 0 {
		for i := 0; i < len(fields)-1; i++ {
			fields[i].HasComma = true
		}
	}

	return fields, indexes
}

func mapField(f *nemgen.Field, dbType config.DBType) SchemaField {
	fieldType := ""
	if dbType == config.MYSQLDBType {
		fieldType = FieldTypeToMYSQL(f)
	}

	notNull := ""
	if f.Required {
		notNull = "NOT NULL"
	}

	ft := SchemaField{
		Name:  f.Identifier,
		Type:  fieldType,
		Field: f,
		Null:  notNull,
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
	return ft
}

func mapFieldsToSelectFields(fields []*nemgen.Field, dbType config.DBType) []SchemaSelectStatementField {
	res := []SchemaSelectStatementField{}
	for _, f := range fields {
		sf := mapField(f, dbType)
		nf := SchemaSelectStatementField{
			Name:   f.Identifier,
			Field:  sf,
			IsLast: false,
		}
		res = append(res, nf)
	}

	if len(res) > 0 {
		res[len(res)-1].IsLast = true
	}

	return res

}
