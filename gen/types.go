package gen

import (
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

type SchemaSelectStatement struct {
	Name             string
	Identifier       string
	EntityIdentifier string
	Fields           []SchemaSelectStatementField
	IsPrimary        bool
	TimeFields       []SchemaField
	SortSupported    bool
}

type SchemaSelectStatementField struct {
	Name   string
	Field  SchemaField
	IsLast bool
}
