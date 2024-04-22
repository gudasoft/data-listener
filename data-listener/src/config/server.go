package config

import (
	"bytes"
	"datalistener/src/logging"
	"datalistener/src/server"
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
	case "parameters":
		return processParametersServerConfig(config.ToMap())
	case "http":
		return processHttpServerConfig(config.ToMap())
	case "https":
		return processHttpsServerConfig(config.ToMap())
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

// Factory method for creating a server config
func processParametersServerConfig(config map[string]interface{}) server.ServerConfig {
	params := server.NewParametersServerConfig()
	if val, ok := config["enabled"].(bool); ok {
		params.Enabled = val
	} else {
		params.Enabled = false
		return params
	}
	if val, ok := config["name"].(string); ok {
		params.NameIsActive = true
		params.Name = val
	}
	if val, ok := config["concurrency"].(int); ok {
		params.ConcurrencyIsActive = true
		params.Concurrency = val
	}
	if val, ok := config["read_buffer_size_kilobyte"].(int); ok {
		params.ReadBufferSizeKilobyteIsActive = true
		params.ReadBufferSizeKilobyte = val
	}
	if val, ok := config["write_buffer_size_kilobyte"].(int); ok {
		params.WriteBufferSizeKilobyteIsActive = true
		params.WriteBufferSizeKilobyte = val
	}
	if val, ok := config["write_timeout_seconds"].(int); ok {
		params.WriteTimeoutSecondsIsActive = true
		params.WriteTimeoutSeconds = val
	}
	if val, ok := config["idle_timeout_seconds"].(int); ok {
		params.IdleTimeoutSecondsIsActive = true
		params.IdleTimeoutSeconds = val
	}
	if val, ok := config["max_conns_per_ip"].(int); ok {
		params.MaxConnsPerIPIsActive = true
		params.MaxConnsPerIP = val
	}
	if val, ok := config["max_requests_per_conn"].(int); ok {
		params.MaxRequestsPerConnIsActive = true
		params.MaxRequestsPerConn = val
	}
	if val, ok := config["max_keep_alive_duration_seconds"].(int); ok {
		params.MaxKeepAliveDurationSecondsIsActive = true
		params.MaxKeepAliveDurationSeconds = val
	}
	if val, ok := config["max_request_body_size_kilobyte"].(int); ok {
		params.MaxRequestBodySizeKilobyteIsActive = true
		params.MaxRequestBodySizeKilobyte = val
	}
	if val, ok := config["disable_keep_alive"].(bool); ok {
		params.DisableKeepAliveIsActive = true
		params.DisableKeepAlive = val
	}
	if val, ok := config["tcp_keep_alive"].(bool); ok {
		params.TCPKeepAliveIsActive = true
		params.TCPKeepAlive = val
	}
	if val, ok := config["reduce_memory_usage"].(bool); ok {
		params.ReduceMemoryUsageIsActive = true
		params.ReduceMemoryUsage = val
	}
	if val, ok := config["get_only"].(bool); ok {
		params.GetOnlyIsActive = true
		params.GetOnly = val
	}
	if val, ok := config["disable_pre_parse_multipart_form"].(bool); ok {
		params.DisablePreParseMultipartFormIsActive = true
		params.DisablePreParseMultipartForm = val
	}
	if val, ok := config["log_all_errors"].(bool); ok {
		params.LogAllErrorsIsActive = true
		params.LogAllErrors = val
	}
	if val, ok := config["secure_error_log_message"].(bool); ok {
		params.SecureErrorLogMessageIsActive = true
		params.SecureErrorLogMessage = val
	}
	if val, ok := config["disable_header_names_normalizing"].(bool); ok {
		params.DisableHeaderNamesNormalizingIsActive = true
		params.DisableHeaderNamesNormalizing = val
	}
	if val, ok := config["sleep_when_concurrency_limits_exceeded"].(int); ok {
		params.SleepWhenConcurrencyLimitsExceededIsActive = true
		params.SleepWhenConcurrencyLimitsExceededSeconds = val
	}
	if val, ok := config["no_default_server_header"].(bool); ok {
		params.NoDefaultServerHeaderIsActive = true
		params.NoDefaultServerHeader = val
	}
	if val, ok := config["no_default_date"].(bool); ok {
		params.NoDefaultDateIsActive = true
		params.NoDefaultDate = val
	}
	if val, ok := config["no_default_content_type"].(bool); ok {
		params.NoDefaultContentTypeIsActive = true
		params.NoDefaultContentType = val
	}
	if val, ok := config["keep_hijacked_conns"].(bool); ok {
		params.KeepHijackedConnsIsActive = true
		params.KeepHijackedConns = val
	}
	if val, ok := config["close_on_shutdown"].(bool); ok {
		params.CloseOnShutdownIsActive = true
		params.CloseOnShutdown = val
	}
	if val, ok := config["stream_request_body"].(bool); ok {
		params.StreamRequestBodyIsActive = true
		params.StreamRequestBody = val
	}

	return params
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

func processHttpsServerConfig(config map[string]interface{}) server.ServerConfig {
	var httpsConf server.HttpsServerConfig
	httpsConf.Protocol = strings.ToLower(config["protocol"].(string))
	httpsConf.Address = config["address"].(string)
	httpsConf.Port = int(config["port"].(int64))
	httpsConf.ServerCertFile = config["server_tls_certificate"].(string)
	httpsConf.ServerKeyFile = config["server_tls_key"].(string)
	httpsConf.UseMTLS = config["use_mtls"].(bool)
	if httpsConf.UseMTLS {
		httpsConf.CACertFile = config["ca_tls_cert"].(string)
	}
	return &httpsConf
}

func processUnixServerConfig(config map[string]interface{}) server.ServerConfig {
	var unixConf server.UnixServerConfig

	unixConf.Protocol = strings.ToLower(config["protocol"].(string))
	unixConf.Address = config["address"].(string)

	return &unixConf
}
