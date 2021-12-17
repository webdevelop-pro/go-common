package logger

type Config struct {
	LogLevel   string `envconfig:"LOG_LEVEL" default:"info"`
	LogConsole bool   `envconfig:"LOG_CONSOLE" default:"false"`
}
