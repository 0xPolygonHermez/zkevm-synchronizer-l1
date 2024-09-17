package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/etherman"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/log"
	storage "github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/storage"
	syncconfig "github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/config"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
)

const (
	// FlagYes is the flag for yes.
	FlagYes = "yes"
	// FlagCfg is the flag for cfg.
	FlagCfg = "cfg"
	// EnvPrefix is the prefix for the environment variables.
	EnvVarPrefix = "ZKEVM_SYNCL1"
)

type Config struct {
	// Configure Log level for all the services, allow also to store the logs in a file
	Log          log.Config        `mapstructure:"Log"`
	SQLDB        storage.Config    `mapstructure:"SQLDB"`
	Synchronizer syncconfig.Config `mapstructure:"Synchronizer"`
	Etherman     etherman.Config   `mapstructure:"Etherman"`
}

// Default parses the default configuration values.
func Default() (*Config, error) {
	var cfg Config
	viper.SetConfigType("toml")

	err := viper.ReadConfig(bytes.NewBuffer([]byte(DefaultValues)))
	if err != nil {
		return nil, err
	}
	err = viper.Unmarshal(&cfg, viper.DecodeHook(mapstructure.TextUnmarshallerHookFunc()))
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

// Load loads the configuration
func Load(ctx *cli.Context) (*Config, error) {
	configFilePath := ctx.String(FlagCfg)
	return LoadFile(configFilePath)
}

// Load loads the configuration
func LoadFileFromString(configFileData string, configType string) (*Config, error) {
	cfg := &Config{}
	err := loadString(cfg, DefaultValues, "toml", true, EnvVarPrefix, nil)
	if err != nil {
		return nil, err
	}
	expectedKeys := viper.AllKeys()
	err = loadString(cfg, configFileData, configType, true, EnvVarPrefix, &expectedKeys)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func SaveConfigToString(cfg Config) (string, error) {
	b, err := json.Marshal(cfg)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Load loads the configuration
func LoadFile(configFilePath string) (*Config, error) {
	_, fileName := filepath.Split(configFilePath)
	fileExtension := strings.TrimPrefix(filepath.Ext(fileName), ".")
	configData, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}
	cfg, err := LoadFileFromString(string(configData), fileExtension)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

// Load loads the configuration
func loadString(cfg *Config, configData string, configType string, allowEnvVars bool, envPrefix string, expectedKeys *[]string) error {
	viper.SetConfigType(configType)
	if allowEnvVars {
		replacer := strings.NewReplacer(".", "_")
		viper.SetEnvKeyReplacer(replacer)
		viper.SetEnvPrefix(envPrefix)
		viper.AutomaticEnv()
	}
	err := viper.ReadConfig(bytes.NewBuffer([]byte(configData)))
	if err != nil {
		return err
	}
	decodeHooks := []viper.DecoderConfigOption{
		// this allows arrays to be decoded from env var separated by ",", example: MY_VAR="value1,value2,value3"
		viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(mapstructure.TextUnmarshallerHookFunc(), mapstructure.StringToSliceHookFunc(","))),
	}

	err = viper.Unmarshal(&cfg, decodeHooks...)
	if err != nil {
		return err
	}
	if expectedKeys != nil {
		configKeys := viper.AllKeys()
		err = checkUnknownFieldsOnFile(configKeys, *expectedKeys)
		if err != nil {
			return err
		}
	}
	return nil
}

func checkUnknownFieldsOnFile(keysOnFile, expectedConfigKeys []string) error {
	for _, key := range keysOnFile {
		if !contains(expectedConfigKeys, key) {
			return fmt.Errorf("unknown field %s on config file", key)
		}
	}
	return nil
}

func contains(keys []string, key string) bool {
	for _, k := range keys {
		if k == key {
			return true
		}
	}
	return false
}
