package middleware

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/webdevelop-pro/go-common/context/keys"
	"github.com/webdevelop-pro/go-common/verser"
	logger "github.com/webdevelop-pro/go-logger"
)

// SetLogger adds logger to context
func SetLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// get request's context
		ctx := c.Request().Context()
		ipAddress, _ := c.Get(IpAddressContextKey).(string)
		identityID, _ := keys.GetCtxValue(ctx, keys.IdentityID).(string)

		logInfo := logger.ServiceContext{
			Service: verser.GetService(),
			Version: verser.GetVersion(),
			SourceReference: &logger.SourceReference{
				Repository: verser.GetRepository(),
				RevisionID: verser.GetRevisionID(),
			},
			User: identityID,
			HttpRequest: &logger.HttpRequestContext{
				Method:    c.Request().Method,
				RemoteIp:  ipAddress,
				URL:       c.Request().RequestURI,
				UserAgent: c.Request().UserAgent(),
				Referrer:  c.Request().Referer(),
			},
		}

		ctx = context.WithValue(ctx, logger.ServiceContextInfo, logInfo)

		// create sub logger
		log := logger.NewComponentLogger("echo", c)
		// add logger to context
		ctx = log.WithContext(ctx)
		// enrich context
		c.SetRequest(c.Request().WithContext(ctx))
		// next handler
		return next(c)
	}
}
