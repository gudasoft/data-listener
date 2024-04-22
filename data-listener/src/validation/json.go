package validation

import "fmt"

type JSONValidationConfig struct {
	Enabled bool
}

func (cfg JSONValidationConfig) Validate() bool {
	return cfg.Enabled
}

func (cfg JSONValidationConfig) String() string {
	return fmt.Sprintf("JSON validation enabled")
}
