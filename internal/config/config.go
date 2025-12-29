package config

import (
	"errors"
	"strings"

	"github.com/go-playground/validator"
	"github.com/spf13/viper"
)

const (
	// general
	ApplicationName = "Arctic"

	// defaults
	DefaultHTTPPort  = 9726
	DefaultLogFormat = "json"
	DefaultLogLevel  = "info"
)

type Config struct {
	Server *ServerConfig `mapstructure:"server" validate:"required"`
	Log    *LogConfig    `mapstructure:"log" validate:"required"`
}

type ServerConfig struct {
	Port int `mapstructure:"port" validate:"required,min=1,max=65535"`
}

type LogConfig struct {
	Level  string `mapstructure:"level" validate:"required,oneof=debug info warn error"`
	Format string `mapstructure:"format" validate:"required,oneof=json text"`
}

// Load loads the static configuration required to start Arctic.
// This includes database connection details, HTTP server settings, logging, and more.
//
// When no config file path is specified, Load searches the following locations:
//
//	/etc/arctic/config.{yml,yaml,json,toml}
//	~/.arctic/config.{yml,yaml,json,toml}
//
// Configuration values are resolved in order of precedence (highest to lowest):
//
//  1. Environment variables (ARCTIC_SERVER_PORT, ARCTIC_LOG_LEVEL, etc.)
//  2. Config file
//  3. Built-in defaults
func Load(configPath string) (*Config, error) {
	vip := viper.New()
	setDefaults(vip)
	envOverride(vip)

	if configPath != "" {
		viper.SetConfigFile(configPath)
	} else {
		viper.SetConfigName("config")
		viper.AddConfigPath("/etc/arctic")
		viper.AddConfigPath("$HOME/.arctic")
	}

	if err := viper.ReadInConfig(); err != nil {
		var fileNotFoundErr viper.ConfigFileNotFoundError
		if !errors.As(err, &fileNotFoundErr) {
			return nil, err
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	v := validator.New()
	if err := v.Struct(cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func setDefaults(_ *viper.Viper) {
	viper.SetDefault("server.port", DefaultHTTPPort)
	viper.SetDefault("log.level", DefaultLogLevel)
	viper.SetDefault("log.format", DefaultLogFormat)
}

func envOverride(_ *viper.Viper) {
	viper.SetEnvPrefix("ARCTIC")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
}
