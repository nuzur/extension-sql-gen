package gen

import (
	"fmt"
	"slices"
	"sort"
	"strings"

	nemgen "github.com/nuzur/extension-sdk/proto_deps/nem/idl/gen"
	"github.com/nuzur/extension-sql-gen/config"
)

type SchemaTemplate struct {
	Entities []SchemaEntity
}

// entity
type SchemaEntity struct {
	DBType           config.DBType
	Name             string
	NameTitle        string
	PrimaryKeys      []string
	Fields           []SchemaField
	Indexes          []SchemaIndex
	Constraints      []SchemaConstraint
	SelectStatements []SchemaSelectStatement
}

func (e SchemaEntity) IsPrimaryKey(fieldIdentifier string) bool {
	return slices.Contains(e.PrimaryKeys, fieldIdentifier)
}

func (e SchemaEntity) PrimaryKeysIdentifiers() string {
	return strings.Join(e.PrimaryKeys, ", ")
}

func (e SchemaEntity) PrimaryKeysWhereClause() string {
	keys := []string{}
	for _, pk := range e.PrimaryKeys {
		if e.DBType == config.MYSQLDBType {
			keys = append(keys, fmt.Sprintf("`%s` = ?", pk))
		} else if e.DBType == config.PGDBType {
			keys = append(keys, fmt.Sprintf(`"%s" = ?`, pk))
		}
	}
	return strings.Join(keys, " AND ")
}

func (e SchemaEntity) UpdateFields() string {
	fields := []string{}
	for _, f := range e.Fields {
		if !slices.Contains(e.PrimaryKeys, f.Name) {
			if e.DBType == config.MYSQLDBType {
				fields = append(fields, fmt.Sprintf("`%s` = ?", f.Name))
			} else if e.DBType == config.PGDBType {
				fields = append(fields, fmt.Sprintf(`"%s" = ?`, f.Name))
			}
		}
	}
	return strings.Join(fields, ", ")
}

// field
type SchemaField struct {
	Name      string
	NameTitle string
	Type      string
	Field     *nemgen.Field
	Null      string
	HasComma  bool
	Default   string
	Unique    string
}

func (f SchemaField) Postfix() string {
	res := []string{}
	if f.Null != "" {
		res = append(res, f.Null)
	}
	if f.Default != "" {
		res = append(res, f.Default)
	}
	return strings.Join(res, " ")
}

// index
type SchemaIndex struct {
	DBType     config.DBType
	Name       string
	FieldNames map[string]string
	Index      *nemgen.Index
	TypePrefix string
	Type       string
	TypeSort   int
	HasComma   bool
}

func (i SchemaIndex) FieldNamesIdentifiers() string {
	fields := i.Index.Fields
	sort.Slice(fields, func(i, j int) bool {
		return fields[i].Priority < fields[j].Priority
	})

	fieldsStr := []string{}
	for _, f := range fields {

		if i.DBType == config.MYSQLDBType {
			order := ""
			if f.Order == nemgen.IndexFieldOrder_INDEX_FIELD_ORDER_ASC {
				order = "ASC"
			} else if f.Order == nemgen.IndexFieldOrder_INDEX_FIELD_ORDER_DESC {
				order = "DESC"
			}
			fieldsStr = append(fieldsStr, fmt.Sprintf("`%s` %s", i.FieldNames[f.FieldUuid], order))
		} else if i.DBType == config.PGDBType {
			fieldsStr = append(fieldsStr, fmt.Sprintf(`"%s"`, i.FieldNames[f.FieldUuid]))
		}
	}

	return fmt.Sprintf("(%s)", strings.Join(fieldsStr, ", "))
}

// select
type SchemaSelectStatement struct {
	Name             string
	Identifier       string
	EntityIdentifier string
	Fields           []SchemaSelectStatementField
	CombinedIndexes  bool
	IsPrimary        bool
	TimeFields       []SchemaField
	SortSupported    bool
}

type SchemaSelectStatementField struct {
	Name   string
	Field  SchemaField
	IsLast bool
}

// contraints
type SchemaConstraint struct {
	DBType       config.DBType
	Name         string
	Relationship *nemgen.Relationship
	TableName    string
	Fields       []SchemaField
	HasComma     bool
}

func (sc SchemaConstraint) ForeignKeyFields() string {
	sort.Slice(sc.Fields, func(i, j int) bool {
		return strings.Compare(sc.Fields[i].Name, sc.Fields[j].Name) < 1
	})
	fields := []string{}
	for _, f := range sc.Fields {
		if sc.DBType == config.MYSQLDBType {
			fields = append(fields, fmt.Sprintf("`%s_%s`", sc.TableName, f.Name))
		} else if sc.DBType == config.PGDBType {
			fields = append(fields, fmt.Sprintf(`"%s_%s"`, sc.TableName, f.Name))
		}
	}

	return strings.Join(fields, ", ")
}

func (sc SchemaConstraint) ReferenceFields() string {
	sort.Slice(sc.Fields, func(i, j int) bool {
		return strings.Compare(sc.Fields[i].Name, sc.Fields[j].Name) < 1
	})
	fields := []string{}
	for _, f := range sc.Fields {
		if sc.DBType == config.MYSQLDBType {
			fields = append(fields, fmt.Sprintf("`%s`", f.Name))
		} else if sc.DBType == config.PGDBType {
			fields = append(fields, fmt.Sprintf(`"%s"`, f.Name))
		}
	}

	return strings.Join(fields, ", ")
}
