# Response component

## Idea
- We want to unify error responses and generic responses accross different microservices.
- Plus, we want to make `ports/http` as independend from `app/app` itself as possible.
- Meaning, `app/app` can also be used in cli utilities, or with grcp handler or graphql.
- So `app/appp` should keep same responses and error codes regardless of the protocol.


At this moment we are heavily using `response.Error` structure:
```go
type Error struct {
	StatusCode int
	Message    map[string][]string
	Err        error
}
```

- StatusCode used to return StatusCode, usually its visible for the end client
- Message - user friendly message
- Err - error structure, which might contain sensetive data or it can be an collection of inheritance of different error messages. Err also contains stack trace if used together with [go-logger](https://github.com/webdevelop-pro/go-logger)


## Usage

Typically, `app/app` will return a `response.Error` structure and handler pass it to an end-client

```go
fileMeta, err := h.app.UploadFile(c, reqData)
if err != nil {
  e := err.(response.Error)
  if e.StatusCode >= 500 {
    log.Error().Stack().Err(err).Msgf("system error happen")
  }

  return c.JSON(e.StatusCode, e.Message)
}
```

Application should create Message using `map[string][]interface{}` type to simplify parsing of the response for the client.

## Consts

Typical error messages:
- `codes.BadRequestMsg` - return `bad request` error to frontend. Used when error message is not accesible or when you don't want to give any detail regarding error to the client
- `codes.NotAuthorizedMsg` - return `authorization error`. Used when used trying to access authorized method without or with incorrect credentials
- `codes.InternalErrMsg` - return `internal/server error`. Used during internal errors when you don't want to give any additional details to the client

Useful methods:
- `response.New` - takes `error`, `statusCode`, `message` as input
```go
err = app.repo.UpdateProfileData(ctx, profile)
if err != nil {
  return response.New(
    err,
    http.StatusInternalServerError,
    response.InternalErrMsg,
  )
}
```

- `response.BadRequest` - create generic `BadRequest` response with custom error handler
```go
if err = c.Bind(&reqData); err != nil {
	return response.BadRequest(err)
}
```

- `response.BadRequestMsg` - create generic `BadRequestMsg` response with client message
```go
if err = c.Decode(&reqData); err != nil {
	return response.BadRequestMsg(fmt.Errorf("cannot decode incoming data"))
}
```

Validation documentetion can be found [here](../validator)

## Error codes convention

- return `http.StatusPreconditionFailed` during validation failer with 3rd party service. Example profile does not have all required field for stripe
- return `http.StatusFailedDependency` during failer with 3rd party service. Example stripe was not able to handle request and return an error
- use `400`, `401`, `402`, `404`, `500` status code based on RFC 9110 documentation

## ToDo

- [ ] for `response.Error` add internal error code and give integer code to every error similar to https://transactapi.readme.io/docs/error-codes
- [ ] Ideally `app/app` should not return `http.StatusXXXX` but return an internal error code and `ports/http` should convert it to http status