package configurator

import (
	"sync"

	"github.com/jinzhu/copier"
)

var copierOption = copier.Option{IgnoreEmpty: true, DeepCopy: true}

// Configurator is a struct for getting/setting a configuration
type Configurator struct {
	mu        sync.Mutex
	configMap map[string]interface{}
}

// NewConfigurator returns a new instance of Configurator
func NewConfigurator(configs ...Configuration) *Configurator {
	c := &Configurator{
		configMap: map[string]interface{}{},
	}

	for _, cfg := range configs {
		c.Set(cfg.Name, cfg.Config)
	}

	return c
}

// Set sets a new configuration to a config map
func (c *Configurator) Set(key string, conf interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.configMap[key] = conf
}

// Get returns or creates a Configuration
func (c *Configurator) Get(key string, conf interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	configuration, exists := c.configMap[key]
	if exists {
		if err := copier.CopyWithOption(conf, configuration, copierOption); err != nil {
			panic(err)
		}

		return
	}

	if err := NewConfiguration(conf); err != nil {
		panic(err)
	}

	c.configMap[key] = conf

	return
}
