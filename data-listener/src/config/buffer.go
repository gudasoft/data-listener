package config

import (
	"bytes"
	dest "datalistener/src/destinations"
	"datalistener/src/destinations/buffer"
	"datalistener/src/logging"
	"fmt"
	"strings"

	"github.com/pelletier/go-toml"
)

func ProcessBufferConfigs(table *toml.Tree) []dest.BufferConfig {
	var configs []dest.BufferConfig
	for _, key := range table.Keys() {
		value := table.Get(key)
		switch config := value.(type) {
		case []*toml.Tree:
			for _, subConfig := range config {
				if enabled, ok := subConfig.Get("enabled").(bool); ok && enabled {
					configs = append(configs, processBufferConfig(key, subConfig))
				}
			}
		case *toml.Tree:
			if enabled, ok := config.Get("enabled").(bool); ok && enabled {
				configs = append(configs, processBufferConfig(key, config))
			}
		default:
			if key != "size_megabyte" && key != "convert_to_jsonl" {
				fmt.Printf("Unhandled buffer configuration type: %s\n", key)
			}
		}
	}
	return configs
}

func processBufferConfig(configType string, config *toml.Tree) dest.BufferConfig {
	switch configType {
	case "file":
		return processFileBufferConfig(config.ToMap())
	case "http":
		return processHttpBufferConfig(config.ToMap())
	case "https":
		return processHttpsBufferConfig(config.ToMap())
	case "kafka":
		return processKafkaBufferConfig(config.ToMap())
	case "sftp":
		return processSftpBufferConfig(config.ToMap())
	case "smtp":
		return processSmtpsBufferConfig(config.ToMap())
	case "s3":
		return processS3BufferConfig(config.ToMap())
	default:
		logging.Logger.Sugar().Debugf("Unhandled buffer configuration type: %s\n", configType)
		return nil
	}
}

func GetBufferConfigInfo(configs []dest.BufferConfig) string {
	var bytesBuffer bytes.Buffer

	for _, config := range configs {
		switch cfg := config.(type) {
		case *buffer.FileBufferConfig:
			bytesBuffer.WriteString(fmt.Sprintf("File writing to: %+v\n", cfg))
		case *buffer.FileBufferConfigUnique:
			bytesBuffer.WriteString(fmt.Sprintf("File writing to: %+v\n", cfg))
		case *buffer.HttpBufferConfig:
			bytesBuffer.WriteString(fmt.Sprintf("HTTP output: %+v\n", cfg))
		case *buffer.HttpsBufferConfig:
			bytesBuffer.WriteString(fmt.Sprintf("HTTPS output: %+v\n", cfg))
		case *buffer.HttpsMtlsBufferConfig:
			bytesBuffer.WriteString(fmt.Sprintf("HTTPS MTLS output: %+v\n", cfg))
		case *buffer.KafkaBufferConfig:
			bytesBuffer.WriteString(fmt.Sprintf("Kafka output: %+v\n", cfg))
		case *buffer.SftpBufferConfig:
			bytesBuffer.WriteString(fmt.Sprintf("SFTP output: %+v\n", cfg))
		case *buffer.SmtpsBufferConfig:
			bytesBuffer.WriteString(fmt.Sprintf("SMTPS output: %+v\n", cfg))
		case *buffer.S3BufferConfig:
			bytesBuffer.WriteString(fmt.Sprintf("S3 output: %+v\n", cfg))
		default:
			bytesBuffer.WriteString(fmt.Sprintf("\nWARNING! Unhandled buffer configuration type: %+v\n", cfg))
		}
	}

	return bytesBuffer.String()
}

func processFileBufferConfig(config map[string]interface{}) dest.BufferConfig {
	if config["unique_file_per_buffer"].(bool) {
		var fileConfig buffer.FileBufferConfigUnique
		fileConfig.UniqueFilePerBuffer = true
		fileConfig.FilePathFormat = config["file_path_format"].(string)
		fileConfig.FileFormat = config["file_format"].(string)
		fileConfig.ItemSeparator = config["item_separator"].(string)
		fileConfig.FileExtension = config["file_extansion"].(string)

		return &fileConfig
	}
	var fileConfig buffer.FileBufferConfig
	fileConfig.UniqueFilePerBuffer = false
	fileConfig.MaxFileSize = int(config["max_file_size_kilobyte"].(int64)) * 1024
	fileConfig.FilePathFormat = config["file_path_format"].(string)
	fileConfig.FileFormat = config["file_format"].(string)
	fileConfig.ItemSeparator = config["item_separator"].(string)
	fileConfig.FileExtension = config["file_extansion"].(string)

	return &fileConfig
}

func processHttpBufferConfig(config map[string]interface{}) dest.BufferConfig {
	var httpConf buffer.HttpBufferConfig
	httpConf.Protocol = strings.ToLower(config["protocol"].(string))
	httpConf.Address = config["address"].(string)
	httpConf.Port = int(config["port"].(int64))
	httpConf.Endpoint = config["endpoint"].(string)
	httpConf.ContentType = config["content_type"].(string)
	httpConf.ItemSeparator = config["item_separator"].(string)

	return &httpConf
}

func processHttpsBufferConfig(config map[string]interface{}) dest.BufferConfig {
	if config["use_mtls"].(bool) {
		var httpsConf buffer.HttpsMtlsBufferConfig
		httpsConf.Protocol = strings.ToLower(config["protocol"].(string))
		httpsConf.Address = config["address"].(string)
		httpsConf.Port = int(config["port"].(int64))
		httpsConf.Endpoint = config["endpoint"].(string)
		httpsConf.ContentType = config["content_type"].(string)
		httpsConf.ItemSeparator = config["item_separator"].(string)
		httpsConf.ClientCertFile = config["client_tls_cert"].(string)
		httpsConf.ClientKeyFile = config["client_tls_key"].(string)
		httpsConf.CACertFile = config["ca_tls_cert"].(string)
		httpsConf.SkipHostNameVerification = config["skip_host_name_verification"].(bool)
		return httpsConf
	}
	var httpsConf buffer.HttpsBufferConfig
	httpsConf.Protocol = strings.ToLower(config["protocol"].(string))
	httpsConf.Address = config["address"].(string)
	httpsConf.Port = int(config["port"].(int64))
	httpsConf.Endpoint = config["endpoint"].(string)
	httpsConf.ContentType = config["content_type"].(string)
	httpsConf.ItemSeparator = config["item_separator"].(string)
	return &httpsConf
}

func processKafkaBufferConfig(config map[string]interface{}) dest.BufferConfig {
	var kafkaConf buffer.KafkaBufferConfig
	kafkaConf.BootstrapServers = toStringSlice(config["bootstrap_servers"].([]interface{}))
	kafkaConf.Topic = config["topic"].(string)

	return &kafkaConf
}

func processSftpBufferConfig(config map[string]interface{}) dest.BufferConfig {
	var sftpConf buffer.SftpBufferConfig
	sftpConf.Host = config["host"].(string)
	sftpConf.Port = int(config["port"].(int64))
	sftpConf.Username = config["username"].(string)
	sftpConf.Password = config["password"].(string)
	sftpConf.RemotePath = config["remote_path"].(string)
	sftpConf.LocalPath = config["local_path"].(string)

	return &sftpConf
}

func processSmtpsBufferConfig(config map[string]interface{}) dest.BufferConfig {
	var smtpsConf buffer.SmtpsBufferConfig
	smtpsConf.Host = config["host"].(string)
	smtpsConf.Port = int(config["port"].(int64))
	smtpsConf.Username = config["username"].(string)
	smtpsConf.Password = config["password"].(string)
	smtpsConf.UseTLS = config["use_tls"].(bool)

	return &smtpsConf
}

func processS3BufferConfig(config map[string]interface{}) dest.BufferConfig {
	var s3Conf buffer.S3BufferConfig
	s3Conf.Region = config["region"].(string)
	s3Conf.Bucket = config["bucket"].(string)
	s3Conf.PrefixFormat = config["prefix_format"].(string)
	s3Conf.KeyFormat = config["key_format"].(string)
	s3Conf.ObjType = config["obj_type"].(string)
	s3Conf.ItemSeparator = config["item_separator"].(string)

	return &s3Conf
}
