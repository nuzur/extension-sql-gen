package gen

type SchemaTemplate struct {
	Entities []SchemaEntity
}

type SchemaEntity struct {
	Name             string
	NameTitle        string
	PrimaryKey       string
	Fields           []SchemaField
	Indexes          []SchemaIndex
	SelectStatements []SchemaSelectStatement
}

type SchemaField struct {
	Name     string
	Type     string
	Null     string
	HasComma bool
	Default  string
	Unique   string
}

type SchemaIndex struct {
	Name      string
	FieldName string
	HasComma  bool
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
