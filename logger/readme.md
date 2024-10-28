# Golang logging enhancements


This project adds enhancements to improve logging in go services. Service based on the [zerolog](https://github.com/rs/zerolog) and [pkg/errors](https://github.com/pkg/errors) with additional improvements:

- add stack trace for the error, panic and fatal errors
- add request context to every logging instance
- hook for seemless integration with [google cloud error reporting](https://cloud.google.com/error-reporting)

## Usage

```go
package main

import (
  "net/http"

  "github.com/pkg/errors"
  "github.com/labstack/echo/v4"
  "github.com/webdevelop-pro/go-common/logger"
)

func errorFunc() error {
  return errors.New("some error")
}

func main() {
  e := echo.New()
  defaultLogger := logger.NewComponentLogger("main", nil) // logger without context
  e.GET("/", func(c echo.Context) error {
    err := errorFunc()
    log := logger.NewComponentLogger("get-func", e) // logger with get request context
    log.Error().Ctx(ctx).Stack().Err(err).Msg("error while getting element with id 123")
    return c.String(http.StatusOK, "Hello, World!")
  })
  defaultLogger.Fatal.Err(e.Start(":1323")).Msg("server went down")
}
```

Will output
```json
{
  "level": "error",
  "component": "get-func",
  "stack": [
    {
      "func": "errorFunc",
      "line": "12",
      "source": "main.go"
    },
    {
      "func": "main.func1",
      "line": "19",
      "source": "main.go"
    },
    {
      "func": "(*Echo).add.func1",
      "line": "575",
      "source": "echo.go"
    },
    {
      "func": "(*Echo).ServeHTTP",
      "line": "662",
      "source": "echo.go"
    },
    {
      "func": "serverHandler.ServeHTTP",
      "line": "2936",
      "source": "server.go"
    },
    {
      "func": "(*conn).serve",
      "line": "1995",
      "source": "server.go"
    },
    {
      "func": "goexit",
      "line": "1172",
      "source": "asm_arm64.s"
    }
  ],
  "error": "user id: 123 not found in user_users",
  "severity": "ERROR",
  "serviceContext": {
    "service": "hello-world",
    "version": "1.2.3-git-sha256",
    "user": "123123",
    "httpRequest": {
      "method": "GET",
      "url": "/",
      "userAgent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36",
      "referrer": "",
      "responseStatusCode": 200,
      "remoteIp": "127.0.0.1:62312"
    }
  },
  "@type": "type.googleapis.com/google.devtools.clouderrorreporting.v1beta1.ReportedErrorEvent",
  "time": "2023-03-20T18:08:42+01:00",
  "message": "row not found"
}
```

## Guidelines
- Always log your error twice: right when it happen with stacktrace (`.Stack()`) and in ports handler (http/pubsub/websocket/etc)

#### Notes
Key error elements:
- `level` and `severity`: error level
- `message`: generic error message, i.e. row now found
- `err`: detail error message, i.e. element 123 not found in database
- `component`: name of the component
- `serviceContext`: service information, including user id and request info


### Config

- `LOG_LEVEL` define log level, required
- `LOG_CONSOLE` add terminal colors, useful for local development. Default false
    
## Contributing
[TBD]

