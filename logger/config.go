package logger

type Config struct {
	LogLevel   string `envconfig:"LOG_LEVEL"`
	LogConsole bool   `envconfig:"LOG_CONSOLE" default:"false"`
}
