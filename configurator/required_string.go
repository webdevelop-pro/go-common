package configurator

import (
	"errors"
	"strings"
)

// RequiredString is a string env value that rejects empty or whitespace-only
// input when decoded by envconfig.
type RequiredString string

// Decode implements envconfig.Decoder.
func (s *RequiredString) Decode(value string) error {
	value = strings.TrimSpace(value)
	if value == "" {
		return errors.New("missing value")
	}
	*s = RequiredString(value)
	return nil
}

func (s RequiredString) String() string {
	return string(s)
}
