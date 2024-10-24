package validator

import (
	"fmt"
	"net/http"

	"github.com/friendsofgo/errors"
	"github.com/go-playground/validator/v10"
	"github.com/webdevelop-pro/go-common/response"
)

const (
	MsgRequired = "missing data for required field"
	MsgEmail    = "not a valid email address"
	MsgMin      = "shorter than minimum length %s"
	MsgMax      = "longer than maximum length %s"
	MsgGte      = "greater than or equal to %s"
	MsgGt       = "greater than %s"
	MsgOneOf    = "must be one of: %s"
	MsgEq       = "must be equal to: %s"
	MsgSSN      = "is a valid social security number: %s"
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
	case "gt":
		return fmt.Sprintf(MsgGt, fe.Param())
	case "gte":
		return fmt.Sprintf(MsgGte, fe.Param())
	case "oneof":
		return fmt.Sprintf(MsgOneOf, fe.Param())
	case "eq":
		return fmt.Sprintf(MsgEq, fe.Param())
	case "ssn":
		return fmt.Sprintf(MsgSSN, fe.Value())
	}
	return fe.Error() // default error
}

// Validate check payloads and return error list
func (va Validator) Verify(i interface{}, httpStatus int) error {
	// call the `Struct` function passing in your payload
	err := va.validator.Struct(i)
	if err != nil {
		fieldErrors := response.Error{
			StatusCode: httpStatus,
			Err:        errors.Wrapf(err, "validator error"),
			Message:    make(map[string][]string),
		}
		var ve validator.ValidationErrors
		strErr := "validator error:"
		if errors.As(err, &ve) {
			for _, fe := range ve {
				fieldName := fe.Field()
				_, ok := fieldErrors.Message[fieldName]
				if ok {
					fieldErrors.Message[fieldName] = append(
						fieldErrors.Message[fieldName],
						beautifulMsg(fe),
					)
				} else {
					fieldErrors.Message[fieldName] = []string{beautifulMsg(fe)}
					strErr = fmt.Sprintf("%s %s %s,", strErr, fieldName, beautifulMsg(fe))
				}
			}
		} else {
			fieldErrors.StatusCode = http.StatusInternalServerError
		}
		strErr = strErr[0 : len(strErr)-1]
		// We change default validator error message cause it does not provide details we need
		fieldErrors.Err = errors.Errorf(strErr)
		return fieldErrors
	}
	return nil
}

// ValidateBadRequest execute Validate with default BadRequest response
func (va Validator) Validate(i interface{}) error {
	return va.Verify(i, http.StatusBadRequest)
}
