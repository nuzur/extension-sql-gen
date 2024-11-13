package gen

import (
	"encoding/json"

	"github.com/nuzur/extension-sql-gen/config"
)

type Metadata struct {
	ConfigValues *config.Values `json:"config_values"`
	DownloadURL  string         `json:"download_url"`
}

func (m Metadata) ToString() string {
	bytes, err := json.Marshal(m)
	if err != nil {
		return "{}"
	}

	return string(bytes)
}
