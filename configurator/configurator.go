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

// New returns a new instance of Configurator
func New(configs ...Configuration) *Configurator {
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
func (c *Configurator) Get(key string, conf interface{}) interface{} {
	c.mu.Lock()
	defer c.mu.Unlock()

	configuration, exists := c.configMap[key]
	if exists {
		if err := copier.CopyWithOption(conf, configuration, copierOption); err != nil {
			panic(err)
		}

		return conf
	}

	if err := NewConfiguration(conf); err != nil {
		panic(err)
	}

	c.configMap[key] = conf

	return conf
}

// New sets a new configuration
func (c *Configurator) New(key string, conf interface{}, prefixes ...string) interface{} {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := NewConfiguration(conf, prefixes...); err != nil {
		panic(err)
	}

	c.configMap[key] = conf

	return conf
}
