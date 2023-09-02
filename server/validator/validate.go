package validator

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type FieldError struct {
	Param   string
	Message []string
}

func beautifulMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "Missing data for required field."
	case "email":
		return "Not a valid email address."
	case "min":
		return fmt.Sprintf("Shorter than minimum length %s.", fe.Param())
	case "max":
		return fmt.Sprintf("Longer than maximum length %s.", fe.Param())
	}
	return fe.Error() // default error
}

type Validator struct {
	validator *validator.Validate
}

// Validate check payloads and return error list
func (va *Validator) Validate(i interface{}) error {
	// call the `Struct` function passing in your payload
	err := va.validator.Struct(i)
	if err != nil {
		fieldErrors := Error{
			StatusCode: http.StatusBadRequest,
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
