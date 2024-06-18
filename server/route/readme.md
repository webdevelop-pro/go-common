# http routing

Define routing handler and create `GetRoutes` function which should return `route.Route`

```golang
type Handler struct {
	app service.App
}

func NewHandler() *Handler {
	h := Handler{}
	return &h
}

func (hg Handler) GetRoutes() []route.Route {
	tokenValidatorMiddleware := middleware.KeyAuth(func(key string, c echo.Context) (bool, error) {
		return key == "valid-key", nil
	})

	return []route.Route{
		{
			Method: http.MethodHead,
			Path:   "/api/v1.0/healthcheck",
			Handle: hg.healthCheck,
		},
		{
			Method:      http.MethodPost,
			Path:        "/v1.0/user",
			Handle:      hg.createUser,
			Middlewares: []echo.MiddlewareFunc{tokenValidatorMiddleware},
		},
  }
}
```
