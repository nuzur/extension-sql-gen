package gen

import (
	"fmt"
	"slices"
	"sort"
	"strings"

	nemgen "github.com/nuzur/extension-sdk/proto_deps/nem/idl/gen"
)

type SchemaTemplate struct {
	Entities []SchemaEntity
}

type SchemaEntity struct {
	Name             string
	NameTitle        string
	PrimaryKeys      []string
	Fields           []SchemaField
	Indexes          []SchemaIndex
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
		keys = append(keys, fmt.Sprintf("`%s` = ?", pk))
	}
	return strings.Join(keys, " AND ")
}

func (e SchemaEntity) UpdateFields() string {
	fields := []string{}
	for _, f := range e.Fields {
		if !slices.Contains(e.PrimaryKeys, f.Name) {
			fields = append(fields, fmt.Sprintf("`%s` = ?", f.Name))
		}
	}
	return strings.Join(fields, ", ")
}

type SchemaField struct {
	Name     string
	Type     string
	Field    *nemgen.Field
	Null     string
	HasComma bool
	Default  string
	Unique   string
}

type SchemaIndex struct {
	Name       string
	FieldNames map[string]string
	Index      *nemgen.Index
	HasComma   bool
}

func (i SchemaIndex) FieldNamesIdentifiers() string {
	fields := i.Index.Fields
	sort.Slice(fields, func(i, j int) bool {
		return fields[i].Priority < fields[j].Priority
	})

	fieldsStr := []string{}
	for _, f := range fields {
		order := ""
		if f.Order == nemgen.IndexFieldOrder_INDEX_FIELD_ORDER_ASC {
			order = "ASC"
		} else if f.Order == nemgen.IndexFieldOrder_INDEX_FIELD_ORDER_DESC {
			order = "DESC"
		}
		fieldsStr = append(fieldsStr, fmt.Sprintf("`%s` %s", i.FieldNames[f.FieldUuid], order))
	}

	return fmt.Sprintf("(%s)", strings.Join(fieldsStr, ", "))
}

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
