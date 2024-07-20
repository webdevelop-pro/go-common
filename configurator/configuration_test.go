package configurator

import (
	"fmt"
	"os"
	"testing"
)

// will load .env automatically
// and check if all required envs are exists
func TestNewConfiguration(t *testing.T) {
	type Config struct {
		Host     string `required:"false" default:"localhost" split_words:"true"`
		Port     uint16 `default:"5432" split_words:"true"`
		User     string `required:"true" split_words:"true"`
		Password string `required:"true" split_words:"true"`
	}

	cfg := &Config{}
	err := NewConfiguration(cfg, "DB")
	if err != nil {
		t.Fatalf("cannot process %s", err.Error())
	}
}

// will load .env automatically
// and check if all required envs are exists
func TestAlternativeEnvPath(t *testing.T) {
	type ServerConfig struct {
		Host string `required:"true"`
		Port uint16 `required:"true"`
	}

	os.Setenv("ENV_FILE", ".env.tests")
	cfg := &ServerConfig{}
	err := NewConfiguration(cfg, "")
	if err != nil {
		t.Fatalf("cannot process %s", err.Error())
	}
}

// will load .env automatically
// and check if all required envs are exists
func TestRequiredEnvs(t *testing.T) {
	type ServerConfig struct {
		MyHost string `required:"true"`
		Port   int16  `required:"true"`
	}
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()

	cfg := &ServerConfig{}
	err := NewConfiguration(cfg, "")
	if err == nil {
		t.Errorf("no env set, should return an error")
	}
}
