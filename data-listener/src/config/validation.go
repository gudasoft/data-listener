package config

import (
	"bytes"
	"datalistener/src/logging"
	"datalistener/src/validation"
	"fmt"

	"github.com/pelletier/go-toml"
)

func ProcessValidationConfigs(table *toml.Tree) []validation.ValidationConfig {
	var configs []validation.ValidationConfig
	for _, key := range table.Keys() {
		value := table.Get(key)
		switch config := value.(type) {
		case []*toml.Tree:
			for _, subConfig := range config {
				if enabled, ok := subConfig.Get("enabled").(bool); ok && enabled {
					configs = append(configs, processValidationConfig(key, subConfig))
				}
			}
		case *toml.Tree:
			if enabled, ok := config.Get("enabled").(bool); ok && enabled {
				configs = append(configs, processValidationConfig(key, config))
			}
		default:
			logging.Logger.Sugar().Fatalf("Unhandled configuration type: %s\n", key)
		}
	}
	return configs
}

func GetValidationConfigInfo(configs []validation.ValidationConfig) string {
	var buffer bytes.Buffer

	for _, config := range configs {
		switch cfg := config.(type) {
		case *validation.JSONValidationConfig:
			buffer.WriteString(fmt.Sprintf("JSON Validation: %+v\n", cfg))
		case *validation.WhitelistConfig:
			buffer.WriteString(fmt.Sprintf("Whitelist: %+v\n", cfg))
		default:
			logging.Logger.Sugar().Debugf("Unhandled validation configuration type: %s\n", cfg)
		}
	}
	return buffer.String()
}

func processValidationConfig(configType string, config *toml.Tree) validation.ValidationConfig {
	switch configType {
	case "json":
		return processJSONValidationConfig(config.ToMap())
	case "whitelist":
		return processWhitelistServerConfig(config.ToMap())
	default:
		logging.Logger.Sugar().Debugf("Unhandled validation configuration type: %s\n", configType)
		return nil
	}
}

func processJSONValidationConfig(config map[string]interface{}) validation.ValidationConfig {
	var cfg validation.JSONValidationConfig
	cfg.Enabled = config["enabled"].(bool)
	return &cfg
}

func processWhitelistServerConfig(config map[string]interface{}) validation.ValidationConfig {
	var cfg validation.WhitelistConfig
	cfg.Enabled = config["enabled"].(bool)
	if networks, ok := config["networks"].([]interface{}); ok {
		var networksStr []string
		for _, network := range networks {
			networksStr = append(networksStr, network.(string))
		}
		cfg.Networks = networksStr
	}
	return &cfg
}
