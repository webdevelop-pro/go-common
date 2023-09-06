package validator

import (
	"reflect"
	"strings"
)

func ParamName(fld reflect.StructField) string {
	name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
	if name == "" {
		name = strings.SplitN(fld.Tag.Get("param"), ",", 2)[0]
	}
	if name == "" {
		name = strings.SplitN(fld.Tag.Get("form"), ",", 2)[0]
	}
	if name == "-" {
		return ""
	}

	return name
}
