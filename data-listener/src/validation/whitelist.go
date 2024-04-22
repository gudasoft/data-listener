package validation

import (
	"fmt"
	"strings"
)

type WhitelistConfig struct {
	Enabled  bool
	Networks []string
}

func (cfg WhitelistConfig) Validate() bool {
	return cfg.Enabled
}

func (cfg WhitelistConfig) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("enabled, "))
	sb.WriteString("networks:\n")

	for i, network := range cfg.Networks {
		if i > 0 && i%3 == 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(network + "  ")
	}

	return sb.String()
}
