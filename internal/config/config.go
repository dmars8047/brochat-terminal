package config

const DEFAULT_CONFIG_DIRECTORY_NAME = ".broterm"
const CONFIG_FILE_NAME = "config.json"

type ConfigSettings struct {
	Theme          string `json:"theme"`
	LoggingEnabled bool   `json:"logging_enabled"`
}

func NewConfigSettings() *ConfigSettings {
	return &ConfigSettings{
		Theme:          "default",
		LoggingEnabled: true,
	}
}
