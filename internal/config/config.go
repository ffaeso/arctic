package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

const (
	// general
	ApplicationName = "Arctic"

	// defaults
	DefaultHttpAddr  = 9726
	DefaultLogFormat = "json"
	DefaultLogLevel  = "info"
)

type Config struct {
	Server     ServerConfig     `mapstructure:"server" validate:"required"`
	Log        LogConfig        `mapstructure:"log" validate:"required"`
	Datasource DatasourceConfig `mapstructure:"datasource" validate:"required"`
}

type ServerConfig struct {
	Addr int `mapstructure:"addr" validate:"required,min=1,max=65535"`
}

type LogConfig struct {
	Level  string `mapstructure:"level" validate:"required,oneof=debug info warn error"`
	Format string `mapstructure:"format" validate:"required,oneof=json text"`
}

type DatasourceConfig struct {
	Username string `mapstructure:"username" validate:"required"`
	Password string `mapstructure:"password" validate:"required"`
	Host     string `mapstructure:"host" validate:"required,hostname|ip"`
	Port     int    `mapstructure:"port" validate:"required,min=1,max=65535"`
	DbName   string `mapstructure:"dbname" validate:"required"`
	Sslmode  string `mapstructure:"sslmode" validate:"required,oneof=disable require verify-ca verify-full"`
}

func (p *DatasourceConfig) DSN() string {
	// dsnExample := "postgres://username:password@localhost:5432/database_name?sslmode=required"
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		p.Username,
		p.Password,
		p.Host,
		p.Port,
		p.DbName,
		p.Sslmode,
	)
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
//  1. Environment variables (ARCTIC_SERVER_ADDR, ARCTIC_LOG_LEVEL, etc.)
//  2. Config file
//  3. Built-in defaults
func Load(configPath string) (*Config, error) {
	vip := viper.New()
	setDefaults(vip)
	envOverride(vip)

	if configPath != "" {
		vip.SetConfigFile(configPath)
	} else {
		vip.SetConfigName("config")
		vip.AddConfigPath("/etc/arctic")
		vip.AddConfigPath("$HOME/.arctic")
	}

	if err := vip.ReadInConfig(); err != nil {
		var fileNotFoundErr viper.ConfigFileNotFoundError
		if !errors.As(err, &fileNotFoundErr) {
			return nil, err
		}
	}

	var cfg Config
	if err := vip.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	v := validator.New()
	if err := v.Struct(cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func setDefaults(v *viper.Viper) {
	// server defaults
	v.SetDefault("server.addr", DefaultHttpAddr)

	// logger defaults
	v.SetDefault("log.level", DefaultLogLevel)
	v.SetDefault("log.format", DefaultLogFormat)

	// datasource defaults - empty values for required fields, reasonable defaults for operational settings
	v.SetDefault("datasource.username", "")
	v.SetDefault("datasource.password", "")
	v.SetDefault("datasource.host", "")
	v.SetDefault("datasource.port", 5432)
	v.SetDefault("datasource.dbname", "")
	v.SetDefault("datasource.sslmode", "require")
}

func envOverride(v *viper.Viper) {
	v.SetEnvPrefix("ARCTIC")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// bind all registered keys to environment variables
	for _, key := range v.AllKeys() {
		v.BindEnv(key)
	}
}
