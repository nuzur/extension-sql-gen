package config

type Action string

const (
	SelectAllAction Action = "select-all"
	InsertAction    Action = "insert"
	updateAction    Action = "update"
	DeleteAction    Action = "delete"
	CreateAction    Action = "create"
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
