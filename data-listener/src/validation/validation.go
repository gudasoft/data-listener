package validation

import (
	"datalistener/src/logging"
	"fmt"
	"reflect"
)

type ValidationConfig interface {
	Validate() bool
}

func GetValidationConfig(validationCongfigs []ValidationConfig) (JSONValidationConfig, WhitelistConfig) {
	var jsonValidationConfig JSONValidationConfig
	var whitelistConfig WhitelistConfig
	for _, config := range validationCongfigs {
		switch c := config.(type) {
		case *JSONValidationConfig:
			jsonValidationConfig = *c
		case *WhitelistConfig:
			whitelistConfig = *c
			fmt.Println(whitelistConfig)
		default:
			logging.Logger.Sugar().Debugf("Unhandled validation configuration type: %s\n", reflect.TypeOf(config))
		}
	}

	return jsonValidationConfig, whitelistConfig
}
