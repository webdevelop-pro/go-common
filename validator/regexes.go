package validator

import "regexp"

var (
	pathRegex = regexp.MustCompile(`^[a-zA-Z0-9/]+$`)
)
