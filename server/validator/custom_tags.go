package validator

import (
	"reflect"
	"strings"
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
