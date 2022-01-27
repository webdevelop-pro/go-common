package configurator

import "github.com/kelseyhightower/envconfig"

// Configuration  is a struct for storing configuration
type Configuration struct {
	Name   string
	Config interface{}
}

// NewConfiguration sets conf from env
func NewConfiguration(conf interface{}, prefixes ...string) error {
	prefix := ""

	if len(prefixes) > 0 {
		prefix = prefixes[0]
	}

	if err := envconfig.Process(prefix, conf); err != nil {
		_ = envconfig.Usage(prefix, conf)

		return err
	}

	return nil
}
