package config

import (
	"buffer-handler/logging"
	"buffer-handler/server"
	"bytes"
	"fmt"
	"strings"

	"github.com/pelletier/go-toml"
)

func ProcessServerConfigs(table *toml.Tree) []server.ServerConfig {
	var configs []server.ServerConfig
	for _, key := range table.Keys() {
		value := table.Get(key)
		switch config := value.(type) {
		case []*toml.Tree:
			for _, subConfig := range config {
				if enabled, ok := subConfig.Get("enabled").(bool); ok && enabled {
					configs = append(configs, processServerConfig(key, subConfig))
				}
			}
		case *toml.Tree:
			if enabled, ok := config.Get("enabled").(bool); ok && enabled {
				configs = append(configs, processServerConfig(key, config))
			}
		default:
			logging.Logger.Sugar().Fatalf("Unhandled configuration type: %s\n", key)
		}
	}
	return configs
}

func processServerConfig(configType string, config *toml.Tree) server.ServerConfig {
	switch configType {
	case "http":
		return processHttpServerConfig(config.ToMap())
	case "https":
		return processHhttpsServerConfig(config.ToMap())
	case "prometheus":
		return processPrometheusServerConfig(config.ToMap())
	case "unixsocket":
		return processUnixServerConfig(config.ToMap())
	default:
		logging.Logger.Sugar().Debugf("Unhandled sever configuration type: %s\n", configType)
		return nil
	}
}

func GetServerConfigInfo(configs []server.ServerConfig) string {
	var buffer bytes.Buffer

	for _, config := range configs {
		switch cfg := config.(type) {
		case *server.HttpServerConfig:

			buffer.WriteString(fmt.Sprintf("Listening on: %+v\n", cfg))
		case *server.HttpsServerConfig:

			buffer.WriteString(fmt.Sprintf("Listening on: %+v\n", cfg))
		case *server.PrometheusServerConfig:

			buffer.WriteString(fmt.Sprintf("Prometheus listening on: %+v\n", cfg))
		}
	}

	return buffer.String()
}

func processPrometheusServerConfig(config map[string]interface{}) server.ServerConfig {
	var promConf server.PrometheusServerConfig
	promConf.Address = config["address"].(string)
	promConf.Port = int(config["port"].(int64))
	promConf.Path = config["path"].(string)
	return &promConf
}

func processHttpServerConfig(config map[string]interface{}) server.ServerConfig {
	var httpConf server.HttpServerConfig

	httpConf.Protocol = strings.ToLower(config["protocol"].(string))
	httpConf.Address = config["address"].(string)
	httpConf.Port = int(config["port"].(int64))

	return &httpConf
}

func processHhttpsServerConfig(config map[string]interface{}) server.ServerConfig {
	var httpsConf server.HttpsServerConfig
	httpsConf.Protocol = strings.ToLower(config["protocol"].(string))
	httpsConf.Address = config["address"].(string)
	httpsConf.Port = int(config["port"].(int64))
	httpsConf.TlsCert = config["tls_cert"].(string)
	httpsConf.TlsKey = config["tls_key"].(string)

	return &httpsConf
}

func processUnixServerConfig(config map[string]interface{}) server.ServerConfig {
	var unixConf server.UnixServerConfig

	unixConf.Protocol = strings.ToLower(config["protocol"].(string))
	unixConf.Address = config["address"].(string)

	return &unixConf
}
