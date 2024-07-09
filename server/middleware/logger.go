package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/webdevelop-pro/go-common/context/keys"
	"github.com/webdevelop-pro/go-common/logger"
	"github.com/webdevelop-pro/go-common/verser"
)

// SetLogger adds logger to context
func SetLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// get request's context
		ctx := c.Request().Context()
		ipAddress, _ := c.Get(keys.IPAddress).(string)
		identityID, _ := keys.GetCtxValue(ctx, keys.IdentityID).(string)
		requestID, _ := keys.GetCtxValue(ctx, keys.RequestID).(string)
		msgID, _ := keys.GetCtxValue(ctx, keys.MSGID).(string)

		logInfo := logger.ServiceContext{
			Service: verser.GetService(),
			Version: verser.GetVersion(),
			SourceReference: &logger.SourceReference{
				Repository: verser.GetRepository(),
				RevisionID: verser.GetRevisionID(),
			},
			User:      identityID,
			RequestID: requestID,
			MSGID:     msgID,
			HTTPRequest: &logger.HTTPRequestContext{
				Method:    c.Request().Method,
				RemoteIP:  ipAddress,
				URL:       c.Request().RequestURI,
				UserAgent: c.Request().UserAgent(),
				Referrer:  c.Request().Referer(),
			},
		}

		ctx = keys.SetCtxValue(ctx, keys.LogInfo, logInfo)

		// create sub logger
		log := logger.NewComponentLogger(c.Request().Context(), "echo")
		// add logger to context
		ctx = log.WithContext(ctx)
		// enrich context
		c.SetRequest(c.Request().WithContext(ctx))
		// next handler
		return next(c)
	}
}
