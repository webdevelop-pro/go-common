package middleware

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"strings"

	"github.com/labstack/echo/v4"
)

type contextKey int

const jwtKey contextKey = 1

// JWTPayload is a struct with token claims
type JWTPayload struct {
	UserID string `json:"sub"`
}

// ParseJWTPayload decodes JWT
func ParseJWTPayload(token string) (pld JWTPayload, err error) {
	tokenParts := strings.Split(token, ".")

	if len(tokenParts) != 3 {
		err = fmt.Errorf("not enough parts in token: %s", token)
		return
	}

	err = json.NewDecoder(
		base64.NewDecoder(base64.RawURLEncoding, strings.NewReader(tokenParts[1])),
	).Decode(&pld)

	return
}

// SetJWTPayload is a function which set jwt payload to context
func SetJWTPayload(c echo.Context, pld JWTPayload) {
	ctx := c.Request().Context()
	ctx = context.WithValue(ctx, jwtKey, pld)

	// set user_id to logger
	l := log.Ctx(ctx).With().Str("user_id", pld.UserID).Logger()
	ctx = l.WithContext(ctx)

	c.SetRequest(c.Request().WithContext(ctx))
}

// GetJWTPayload is a function which  extract jwt payload from context
func GetJWTPayload(ctx context.Context) JWTPayload {
	pld, ok := ctx.Value(jwtKey).(JWTPayload)
	if !ok {
		return JWTPayload{}
	}
	return pld
}

// ExtractTokenFromString from string
func ExtractTokenFromString(headerAuth string) string {
	header := strings.Split(headerAuth, " ")

	if len(header) == 2 {
		return header[1]
	}

	return header[0]
}
