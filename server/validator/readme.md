# Response component

## Idea
This is a wrapper around go-playground [validator](https://github.com/go-playground/validator) library.
What changes we providing:
- Simplify error structures to return `{field: [err1, err2]}` and simplify error parcing
- Short and simple error messages to increase error readability for end users

## Usage

- Bind echo framework with validator, `err` is a `response.Error` 
```go
import "github.com/webdevelop-pro/go-common/server/validator"

// get an instance of a validator
e.Validator = validator.New()
err := c.Validate(&requestData)
```

- Or just used it directly in our code:
```go
validate := validator.New()
if err := validate.Validate(dto.ProfileForStripe(profile.Data), http.StatusPreconditionFailed); err != nil {
		return err
}
```
In this example, we are using custom dto structure with requirements for the stripe profile. And in case of the error we creating `response.Error` object with `PreconditionFailed` status


## ToDo
- [ ] add object references in json response. I.e. `{"profile": {"name": ["msg"]}}` instead of `{"name": ["msg"]}`
- [ ] create tag to validate file size
- [ ] create tag to validate file mimetype