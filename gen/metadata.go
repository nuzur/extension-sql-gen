package gen

import (
	"encoding/json"

	"github.com/nuzur/extension-sql-gen/config"
)

type Metadata struct {
	ConfigValues *config.Values `json:"config-values"`
}

func (m Metadata) ToString() string {
	bytes, err := json.Marshal(m)
	if err != nil {
		return "{}"
	}

	return string(bytes)
}
