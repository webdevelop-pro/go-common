package server

/*
func TestHTTPCtx(t *testing.T) {
	ctx := context.Background()
	headers := map[string][]string{
		"X-Request-Id":    {"ZXCasdf123"},
		"X-Forwarded-For": {"31.6.1.12"},
	}

	ctx = keys.SetDefaultHTTPCtx(ctx, headers)

	assert.Equal(t, headers["X-Request-Id"][0], keys.GetCtxValue(ctx, keys.RequestID))
	assert.Equal(t, headers["X-Forwarded-For"][0], keys.GetCtxValue(ctx, keys.IPAddress))
}
*/

/*
ToDo:
	- make actuall request

// If Request Id header is empty we should automatically generate it
func TestEmptyRequestID(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	echoCtx := e.NewContext(req, rec)
	echoCtx.Set(echo.HeaderXRequestID, "123123123")
	// Add middleware
	assert.Equal(t, len(keys.GetCtxValue(echoCtx, keys.RequestID).(string)), 9)
}
*/
