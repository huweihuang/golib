package config

import (
	"fmt"

	log "github.com/huweihuang/golib/logger/zap"
	"github.com/spf13/viper"
)

const (
	defaultConfigPath = "configs"
	defaultConfigType = "yaml"
)

var (
	Get                     = viper.Get
	GetBool                 = viper.GetBool
	GetDuration             = viper.GetDuration
	GetFloat64              = viper.GetFloat64
	GetInt                  = viper.GetInt
	GetInt32                = viper.GetInt32
	GetInt64                = viper.GetInt64
	GetSizeInBytes          = viper.GetSizeInBytes
	GetString               = viper.GetString
	GetStringMap            = viper.GetStringMap
	GetStringMapString      = viper.GetStringMapString
	GetStringMapStringSlice = viper.GetStringMapStringSlice
	GetStringSlice          = viper.GetStringSlice
	GetTime                 = viper.GetTime
	IsSet                   = viper.IsSet
	AllSettings             = viper.AllSettings
	Unmarshal               = viper.Unmarshal
	UnmarshalKey            = viper.UnmarshalKey
)

func Init(configName string) error {
	return InitConfig(defaultConfigPath, configName, defaultConfigType)
}

func InitConfig(configPath, configName, configType string) error {
	viper.SetConfigName(configName)
	viper.SetConfigType(configType)
	viper.AddConfigPath(configPath)
	return viper.ReadInConfig()
}

func InitConfigByPath(configPath string) error {
	viper.SetConfigFile(configPath)
	return viper.ReadInConfig()
}

func InitConfigObject(configName string, configObject interface{}) error {
	filePath := fmt.Sprintf("%s/%s.%s", defaultConfigPath, configName, defaultConfigType)
	return InitConfigObjectByPath(filePath, configObject)
}

func InitConfigObjectByPath(configPath string, configObject interface{}) error {
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read in config by viper, err: %v", err)
	}

	if configObject != nil {
		err := viper.Unmarshal(configObject)
		if err != nil {
			return fmt.Errorf("failed to unmarshal, err: %v", err)
		}
		log.Logger().With("config", configObject).Debug("init config")
	}
	return nil
}
