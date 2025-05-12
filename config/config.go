package config

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
)

func Save() (err error) {
	err = viper.SafeWriteConfig()
	if err != nil {
		var configFileAlreadyExistsError viper.ConfigFileAlreadyExistsError
		if errors.As(err, &configFileAlreadyExistsError) {
			err = viper.WriteConfig()
			if err != nil {
				return fmt.Errorf("failed to overwrite config: %v", err)
			}
		}
	}
	return nil
}
