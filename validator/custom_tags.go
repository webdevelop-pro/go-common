package validator

import (
	"reflect"
	"strings"

	valid "github.com/go-playground/validator/v10"
)

func ParamName(fld reflect.StructField) string {
	two := 2
	name := strings.SplitN(fld.Tag.Get("json"), ",", two)[0]
	if name == "" {
		name = strings.SplitN(fld.Tag.Get("param"), ",", two)[0]
	}
	if name == "" {
		name = strings.SplitN(fld.Tag.Get("form"), ",", two)[0]
	}
	if name == "-" {
		return ""
	}

	return name
}

func isPath(fl valid.FieldLevel) bool {
	return pathRegex.MatchString(fl.Field().String())
}
