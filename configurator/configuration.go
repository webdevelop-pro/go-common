package configurator

import "github.com/kelseyhightower/envconfig"

// Configuration  is a struct for storing configuration
type Configuration struct {
	Name   string
	Config interface{}
}

// NewConfiguration sets conf from env
func NewConfiguration(conf interface{}) error {
	err := envconfig.Process("", conf)
	if err != nil {
		_ = envconfig.Usage("", conf)
		return err
	}
	return nil
}
