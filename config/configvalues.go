package config

type Action string

const (
	SelectSimpleAction             Action = "select-simple"
	SelectForIndexedSimpleAction   Action = "select-indexed-simple"
	SelectForIndexedCombinedAction Action = "select-indexed-combined"
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
