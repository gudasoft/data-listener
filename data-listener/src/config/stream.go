package config

import (
	"bytes"
	dest "datalistener/src/destinations"
	"datalistener/src/destinations/stream"
	"datalistener/src/logging"
	"fmt"
	"strings"

	"github.com/pelletier/go-toml"
)

func ProcessStreamConfigs(table *toml.Tree) []dest.StreamConfig {
	var configs []dest.StreamConfig
	for _, key := range table.Keys() {
		value := table.Get(key)
		switch config := value.(type) {
		case []*toml.Tree:
			for _, subConfig := range config {
				if enabled, ok := subConfig.Get("enabled").(bool); ok && enabled {
					configs = append(configs, processStreamConfig(key, subConfig))
				}
			}
		case *toml.Tree:
			if enabled, ok := config.Get("enabled").(bool); ok && enabled {
				configs = append(configs, processStreamConfig(key, config))
			}
		default:
			logging.Logger.Sugar().Fatalf("Unhandled configuration type: %s\n", key)
		}
	}
	return configs
}

func processStreamConfig(configType string, config *toml.Tree) dest.StreamConfig {
	switch configType {
	case "file":
		return processFileStreamConfig(config.ToMap())
	case "http":
		return processHttpStreamConfig(config.ToMap())
	case "https":
		return processHttpsStreamConfig(config.ToMap())
	case "kafka":
		return processKafkaStreamConfig(config.ToMap())
	case "sftp":
		return processSftpStreamConfig(config.ToMap())
	case "smtp":
		return processSmtpsStreamConfig(config.ToMap())
	case "s3":
		return processS3StreamConfig(config.ToMap())
	default:
		logging.Logger.Sugar().Fatalf("Unhandled stream configuration type: %s\n", configType)
		return nil
	}
}

func GetStreamConfigInfo(configs []dest.StreamConfig) string {
	var buffer bytes.Buffer

	for _, config := range configs {
		switch cfg := config.(type) {
		case *stream.FileStreamConfig:
			buffer.WriteString(fmt.Sprintf("File writing to: %+v\n", cfg))
		case *stream.FileStreamConfigUnique:
			buffer.WriteString(fmt.Sprintf("File writing to: %+v\n", cfg))
		case *stream.HttpStreamConfig:
			buffer.WriteString(fmt.Sprintf("HTTP output: %+v\n", cfg))
		case *stream.HttpsStreamConfig:
			buffer.WriteString(fmt.Sprintf("HTTPS output: %+v\n", cfg))
		case *stream.HttpsMtlsStreamConfig:
			buffer.WriteString(fmt.Sprintf("HTTPS MTLS output: %+v\n", cfg))
		case *stream.KafkaStreamConfig:
			buffer.WriteString(fmt.Sprintf("Kafka output: %+v\n", cfg))
		case *stream.SftpStreamConfig:
			buffer.WriteString(fmt.Sprintf("SFTP output: %+v\n", cfg))
		case *stream.SmtpsStreamConfig:
			buffer.WriteString(fmt.Sprintf("SMTPS output: %+v\n", cfg))
		case *stream.S3StreamConfig:
			buffer.WriteString(fmt.Sprintf("S3 output: %+v\n", cfg))
		default:
			buffer.WriteString(fmt.Sprintf("\nWARNING! Unhandled stream configuration type: %+v\n", cfg))
		}
	}

	return buffer.String()
}

func processFileStreamConfig(config map[string]interface{}) dest.StreamConfig {
	if config["unique_file_per_request"].(bool) {
		var fileConfig stream.FileStreamConfigUnique
		fileConfig.UniqueFilePerRequest = true
		fileConfig.FilePathFormat = config["file_path_format"].(string)
		fileConfig.FileFormat = config["file_format"].(string)
		fileConfig.ItemSeparator = config["item_separator"].(string)
		fileConfig.FileExtension = config["file_extansion"].(string)

		return &fileConfig
	}
	var fileConfig stream.FileStreamConfig
	fileConfig.UniqueFilePerRequest = false
	fileConfig.MaxFileSize = int(config["max_file_size_kilobyte"].(int64)) * 1024
	fileConfig.FilePathFormat = config["file_path_format"].(string)
	fileConfig.FileFormat = config["file_format"].(string)
	fileConfig.ItemSeparator = config["item_separator"].(string)
	fileConfig.FileExtension = config["file_extansion"].(string)

	return &fileConfig
}

func processHttpStreamConfig(config map[string]interface{}) dest.StreamConfig {
	var httpConf stream.HttpStreamConfig
	httpConf.Protocol = strings.ToLower(config["protocol"].(string))
	httpConf.Address = config["address"].(string)
	httpConf.Port = int(config["port"].(int64))
	httpConf.Endpoint = config["endpoint"].(string)
	httpConf.ContentType = config["content_type"].(string)

	return &httpConf
}

func processHttpsStreamConfig(config map[string]interface{}) dest.StreamConfig {
	if config["use_mtls"].(bool) {
		fmt.Println("HTTPS MTLS Stream")
		var httpsConf stream.HttpsMtlsStreamConfig
		httpsConf.Protocol = strings.ToLower(config["protocol"].(string))
		httpsConf.Address = config["address"].(string)
		httpsConf.Port = int(config["port"].(int64))
		httpsConf.EndPoint = config["endpoint"].(string)
		httpsConf.ContentType = config["content_type"].(string)
		httpsConf.ClientCertFile = config["client_tls_cert"].(string)
		httpsConf.ClientKeyFile = config["client_tls_key"].(string)
		httpsConf.CACertFile = config["ca_tls_cert"].(string)
		httpsConf.SkipHostNameVerification = config["skip_host_name_verification"].(bool)
		return &httpsConf
	}
	var httpsConf stream.HttpsStreamConfig
	httpsConf.Protocol = strings.ToLower(config["protocol"].(string))
	httpsConf.Address = config["address"].(string)
	httpsConf.Port = int(config["port"].(int64))
	httpsConf.Endpoint = config["endpoint"].(string)
	httpsConf.ContentType = config["content_type"].(string)
	return &httpsConf
}

func processKafkaStreamConfig(config map[string]interface{}) dest.StreamConfig {
	var kafkaConf stream.KafkaStreamConfig
	kafkaConf.BootstrapServers = toStringSlice(config["bootstrap_servers"].([]interface{}))
	kafkaConf.Topic = config["topic"].(string)

	return &kafkaConf
}

func processSftpStreamConfig(config map[string]interface{}) dest.StreamConfig {
	var sftpConf stream.SftpStreamConfig
	sftpConf.Host = config["host"].(string)
	sftpConf.Port = int(config["port"].(int64))
	sftpConf.Username = config["username"].(string)
	sftpConf.Password = config["password"].(string)
	sftpConf.RemotePath = config["remote_path"].(string)
	sftpConf.LocalPath = config["local_path"].(string)

	return &sftpConf
}

func processSmtpsStreamConfig(config map[string]interface{}) dest.StreamConfig {
	var smtpsConf stream.SmtpsStreamConfig
	smtpsConf.Host = config["host"].(string)
	smtpsConf.Port = int(config["port"].(int64))
	smtpsConf.Username = config["username"].(string)
	smtpsConf.Password = config["password"].(string)
	smtpsConf.UseTLS = config["use_tls"].(bool)

	return &smtpsConf
}

func processS3StreamConfig(config map[string]interface{}) dest.StreamConfig {
	var s3Conf stream.S3StreamConfig
	s3Conf.Region = config["region"].(string)
	s3Conf.Bucket = config["bucket"].(string)
	s3Conf.PrefixFormat = config["prefix_format"].(string)
	s3Conf.KeyFormat = config["key_format"].(string)
	s3Conf.ObjType = config["obj_type"].(string)

	return &s3Conf
}
