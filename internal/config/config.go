package config

import (
	"fmt"

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

	if configPath != "" {
		viper.SetConfigFile(configPath)
	} else {
		// look for 'config.yaml'
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")

		// look for 'config.yaml' in the following locations
		viper.AddConfigPath(".")
		viper.AddConfigPath("$HOME/.arctic")
		viper.AddConfigPath("/etc/arctic/")
	}

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
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

