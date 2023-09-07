package validator

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"github.com/go-playground/validator/v10"
	"github.com/webdevelop-pro/go-common/server/response"
)

const (
	MsgRequired = "Missing data for required field."
	MsgEmail    = "Not a valid email address."
	MsgMin      = "Shorter than minimum length %s."
	MsgMax      = "Longer than maximum length %s."
	MsgGte      = "Greater than or equal to %s."
	MsgOneOf    = "Must be one of: %s."
)

type FieldError struct {
	Param   string
	Message []string
}

type Validator struct {
	validator *validator.Validate
}

func New() Validator {

	v := validator.New()
	v.RegisterTagNameFunc(ParamName)

	return Validator{
		validator: v,
	}
}

func beautifulMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return MsgRequired
	case "email":
		return MsgEmail
	case "min":
		return fmt.Sprintf(MsgMin, fe.Param())
	case "max":
		return fmt.Sprintf(MsgMax, fe.Param())
	case "gte":
		return fmt.Sprintf(MsgGte, fe.Param())
	case "oneof":
		return fmt.Sprintf(MsgOneOf, fe.Param())
	}
	return fe.Error() // default error
}

// Validate check payloads and return error list
func (va Validator) Validate(i interface{}) error {
	// call the `Struct` function passing in your payload
	err := va.validator.Struct(i)
	if err != nil {
		fieldErrors := response.Error{
			StatusCode: http.StatusBadRequest,
			Err:        errors.Wrapf(err, "validator error"),
			Message:    make(map[string][]string),
		}
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			for _, fe := range ve {
				msg, ok := fieldErrors.Message[fe.Field()]
				if ok {
					msg = append(msg, beautifulMsg(fe))
				} else {
					fieldErrors.Message[fe.Field()] = []string{beautifulMsg(fe)}
				}
			}
		} else {
			fieldErrors.StatusCode = http.StatusInternalServerError
		}
		return fieldErrors
	}
	return nil
}
