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
```json
  // MsgBadRequest used to indicate error in incoming data
	MsgBadRequest = map[string][]string{"__error__": {"Bad request."}}
	// MsgNotFound typically used when element haven't been found
	MsgNotFound = map[string][]string{"__error__": {"Not found."}}
	// MsgUnauthorized signalizes lack of token or other authorization data
	MsgUnauthorized = map[string][]string{"__error__": {"Authorization error."}}
	// MsgForbidden used when user does not have any permissions to perform action
	MsgForbidden = map[string][]string{"__error__": {"Forbidden error."}}
	// MsgInternalErr server side error
	MsgInternalErr = map[string][]string{"__error__": {"Internal/server error."}}
```

Useful methods:
- `response.New` - takes `error`, `statusCode`, `message` as input
```go
err = app.repo.UpdateProfileData(ctx, profile)
if err != nil {
  return response.New(
    err,
    http.StatusInternalServerError,
    response.MsgInternalErr,
  )
}
```

- `response.BadRequest` - create generic `BadRequest` response with custom error handler
```go
if err = c.Bind(&reqData); err != nil {
	return response.BadRequest(err, "")
}
```
or 
```go
if err = c.Decode(&reqData); err != nil {
	return response.BadRequest(nil, fmt.Errorf("cannot decode incoming data"))
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