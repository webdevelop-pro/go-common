package response

import (
	"encoding/json"
)

const defaultError = "Something went wrong"

// Error is generic error struct
type Error struct {
	Code        string `json:"code,omitempty"`
	Error       string `json:"error,omitempty"`
	TraceID     string `json:"trace_id,omitempty"`
	Description string `json:"description,omitempty"`
}

//MarshalJSON is custom Marshal function
func (e Error) MarshalJSON() ([]byte, error) {
	type Alias Error

	if e.Error == "" {
		e.Error = defaultError
	}

	return json.Marshal(&struct {
		Alias
	}{
		(Alias)(e),
	})
}
