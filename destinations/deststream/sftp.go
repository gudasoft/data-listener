package deststream

import (
	"buffer-handler/logging"
	"buffer-handler/models"
)

type SftpStreamConfig struct {
	Host       string
	Port       int
	Username   string
	Password   string
	RemotePath string
	LocalPath  string
}

func (cfg SftpStreamConfig) Notify(models.EntryData) error {
	logging.Logger.Debug("Sftp Streaming")

	return nil
}
