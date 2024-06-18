# Log 3rd Party requests

We save all 3rd party requests in [log_logs](https://github.com/webdevelop-pro/i-migration-service/blob/dev/migrations/04_log_logs/01_init.sql#L16) table, and log them to stdout (using [go-logger](https://github.com/webdevelop-pro/go-common/logger) lib)

## Usage

Add Invoke to your fx.Start declaration in main function

```
import (
  client_middleware "github.com/webdevelop-pro/go-common/client/middleware"
)

fx.New(
  ...
  fx.Invoke(
    ...
    func(dwolla dwolla.DwollaWrapper, db client_middleware.DB) {
        dwolla.SetHttpClient(client_middleware.CreateHttpClient("dwolla", db))
    },
  ),
)

```

## TODO
- [ ] Use the same middleware for log incoming and outcomming requests
- [ ] Found way how to deal with content_type_id and object_id

## Diagram
![image](https://github.com/webdevelop-pro/go-common/assets/10445445/09295949-e76f-4303-8c44-c45f699ae266)
