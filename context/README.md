# context

This package contains keys for working with context

For extract some value from context use:

```golang
import (
	"github.com/webdevelop-pro/go-common/context/keys"
    ...
)

...
requestID := keys.GetAsString(ctx, keys.RequestID)
// or
requestID := keys.GetCtxValue(ctx, keys.LogObjectID).(int)

```