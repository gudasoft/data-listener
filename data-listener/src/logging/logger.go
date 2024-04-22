package logging

import (
	"fmt"
	"path/filepath"

	"go.uber.org/zap"
)

type LogConfig struct {
	Path          string
	FileName      string
	Mode          string
	SystemLogging bool
}

var (
	Logger *zap.Logger
)

func (cfg *LogConfig) InitLogger() (*zap.Logger, error) {

	if cfg.Mode == "Disabled" || cfg.Mode == "" {
		cfg.Mode = "Disabled"
		Logger = zap.NewNop()
		return Logger, nil
	}

	var zapConfig zap.Config

	switch cfg.Mode {
	case "Development":
		zapConfig = zap.NewDevelopmentConfig()

	case "Production":
		zapConfig = zap.NewProductionConfig()

	default:
		return nil, fmt.Errorf("Unsupported logger mode: %s", cfg.Mode)
	}

	logFilePath := filepath.Join(cfg.Path, cfg.FileName)
	if cfg.SystemLogging {
		zapConfig.OutputPaths = append(zapConfig.OutputPaths, logFilePath)
	} else {
		zapConfig.OutputPaths = []string{logFilePath}
	}
	//zapConfig.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)

	var err error
	Logger, err = zapConfig.Build()
	if err != nil {
		return nil, err
	}

	return Logger, nil
}

func (cfg LogConfig) String() string {
	if cfg.Mode == "Disabled" {
		return fmt.Sprintln("Logging disabled.")
	}

	str := fmt.Sprintf("Logging at "+
		"%s/%s"+
		", %s mode",
		cfg.Path, cfg.FileName, cfg.Mode)

	if cfg.SystemLogging {
		str += ", printing to console."
		if cfg.Mode == "Development" {
			str += "\nWarning! Development mode and printing to console have perfomance hit. Adjust the settings in the config.toml"
		}
	} else {
		str += ""
	}
	return str
}
