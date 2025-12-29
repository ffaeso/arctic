package logger

type Config struct {
	Level  string `mapstructure:"level" validate:"required"`
	Format string `mapstructure:"format" validate:"required"`
}
