package config

import (
	dest "buffer-handler/destinations"
	"buffer-handler/logging"
	"buffer-handler/server"

	"github.com/pelletier/go-toml"
)

func LoadConfigs(configFile string, logConfig *logging.LogConfig, serverConfigs *[]server.ServerConfig, streamerConfigs *[]dest.StreamConfig, buffererConfigs *[]dest.BufferConfig, bufferSize *int, convertToJSONL *bool) {
	cfg, err := toml.LoadFile(configFile)
	if err != nil {
		logging.Logger.Sugar().Fatalf("Error loading config: %s", err)
		return
	}

	logTable := cfg.Get("logger").(*toml.Tree)
	*logConfig = processLogConfig(logTable)

	serverTable := cfg.Get("server").(*toml.Tree)
	*serverConfigs = ProcessServerConfigs(serverTable)

	streamerTable := cfg.Get("streamer").(*toml.Tree)
	*streamerConfigs = ProcessStreamConfigs(streamerTable)

	buffererTable := cfg.Get("bufferer").(*toml.Tree)
	*buffererConfigs = ProcessBufferConfigs(buffererTable)

	*bufferSize = int(cfg.Get("bufferer.size_megabyte").(int64) * 1024 * 1024)
	*convertToJSONL = bool(cfg.Get("bufferer.convert_to_jsonl").(bool))
}

func processLogConfig(table *toml.Tree) logging.LogConfig {
	var logConfig logging.LogConfig
	logConfig.Path = table.Get("path").(string)
	logConfig.FileName = table.Get("file_name").(string)
	logConfig.Mode = table.Get("mode").(string)
	logConfig.SystemLogging = table.Get("printing_logs_to_console").(bool)
	return logConfig
}

func toStringSlice(interfaces []interface{}) []string {
	result := make([]string, len(interfaces))
	for i, v := range interfaces {
		result[i] = v.(string)
	}
	return result
}
