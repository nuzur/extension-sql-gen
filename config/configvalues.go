package config

type Action string

const (
	SelectSimpleAction             Action = "select_simple"
	SelectForIndexedSimpleAction   Action = "select_indexed_simple"
	SelectForIndexedCombinedAction Action = "select_indexed_combined"
	InsertAction                   Action = "insert"
	UpdateAction                   Action = "update"
	DeleteAction                   Action = "delete"
	CreateAction                   Action = "create"
)

type DBType string

const (
	MYSQLDBType DBType = "mysql"
	PGDBType    DBType = "pg"
)

type Values struct {
	DBType   DBType   `json:"db_type"`
	Entities []string `json:"entities"`
	Actions  []Action `json:"actions"`
}
