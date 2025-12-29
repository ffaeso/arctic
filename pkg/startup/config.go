package config

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server" validate:"required"`
	Postgres PostgresConfig `mapstructure:"postgres" validate:"required"`
	Log      LogConfig      `mapstructure:"log" validate:"required"`
}

type ServerConfig struct {
	Port int `mapstructure:"port" validate:"required"`
}

type PostgresConfig struct {
	Host     string `mapstructure:"host" validate:"required"`
	Port     int    `mapstructure:"port" validate:"required"`
	User     string `mapstructure:"user" validate:"required"`
	Password string `mapstructure:"password" validate:"required"`
	Name     string `mapstructure:"name" validate:"required"`
	SSLMode  string `mapstructure:"sslmode" validate:"required"`
}

type LogConfig struct {
	Level  string `mapstructure:"level" validate:"required"`
	Format string `mapstructure:"format" validate:"required"`
}

func Load(configPath string) (*Config, error) {
	var cfg Config

	// set defaults
	// exclude postgres host, user, password, and dbname
	viper.SetDefault("server.port", 9726)
	viper.SetDefault("postgres.port", 5432)
	viper.SetDefault("postgres.sslmode", "disable")
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "json")

	if configPath != "" {
		viper.SetConfigFile(configPath)
	} else {
		// look for 'config' file
		viper.SetConfigName("config")

		// look for config files in the following locations
		// config can be in json, yml or toml
		// however our app provides yml format as example
		viper.AddConfigPath("$HOME/.arctic") // ~/.arctic/config
		viper.AddConfigPath("/etc/arctic/")  // /etc/arctic/config
	}

	// allow environmental variables to override config file
	//
	// examples:
	// server.port -> ARCTIC_SERVER_PORT
	// postgres.sslmode -> ARCTIC_POSTGRES_SSLMODE
	viper.SetEnvPrefix("ARCTIC")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// these MUST be provided either through config file or environment variables
	viper.BindEnv("postgres.host")     // ARCTIC_POSTGRES_HOST
	viper.BindEnv("postgres.user")     // ARCTIC_POSTGRES_USER
	viper.BindEnv("postgres.password") // ARCTIC_POSTGRES_PASSWORD
	viper.BindEnv("postgres.name")     // ARCTIC_POSTGRES_NAME

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	v := validator.New()
	if err := v.Struct(&cfg); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &cfg, nil
}

func (p *PostgresConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		p.User, p.Password, p.Host, p.Port, p.Name, p.SSLMode,
	)
}
