package configurator

import (
	"testing"
)

type Config struct {
	Host     string `required:"false" default:"localhost" split_words:"true"`
	Port     uint16 `default:"5432" split_words:"true"`
	User     string `required:"true" split_words:"true"`
	Password string `required:"true" split_words:"true"`
}

// will load .env automatically
// and check if all required envs are exists
func TestNewConfiguration(t *testing.T) {
	cfg := &Config{}
	err := NewConfiguration(cfg, "DB")
	if err != nil {
		t.Fatalf("cannot process %s", err.Error())
	}
}
