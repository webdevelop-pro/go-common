package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/webdevelop-pro/go-common/validator"
)

type User struct {
	Email string `json:"email" validate:"required,gt=2"`
}

type handler struct{}

func (h *handler) createUser(e echo.Context) error {
	req := new(User)
	if err := e.Bind(req); err != nil {
		return ErrorBadReqestResponse(e, err)
	}

	err := e.Validate(req)
	if err != nil {
		return ErrorResponse(e, err)
	}
	return e.JSON(http.StatusCreated, req)
}

func TestErrorHandler(t *testing.T) {
	// Setup
	e := echo.New()
	// get an instance of a validator
	e.Validator = validator.New()
	req := httptest.NewRequest(http.MethodPost, "/user", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &handler{}

	res := h.createUser(c)
	// Assertions
	if assert.NoError(t, res) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, "{\"email\":[\"missing data for required field\"]}\n", rec.Body.String())
	}
}
