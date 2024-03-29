package configurator

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// Configuration  is a struct for storing configuration
type Configuration struct {
	Name   string
	Config interface{}
}

// NewConfiguration sets conf from env
func NewConfiguration(conf interface{}, prefixes ...string) error {
	prefix := ""

	err := loadDotEnv()
	if err != nil {
		return err
	}

	if len(prefixes) > 0 {
		prefix = prefixes[0]
	}

	if err := envconfig.Process(prefix, conf); err != nil {
		_ = envconfig.Usage(prefix, conf)

		return err
	}

	return nil
}

func loadDotEnv() error {
	envPath := os.Getenv("ENV_FILE")

	var err error
	if envPath == "" {
		_ = godotenv.Load(".env") // ignore error by default
	} else {
		err = godotenv.Load(envPath) // if path to env file defined, check error
	}

	if err != nil {
		return err
	}

	return nil
}
