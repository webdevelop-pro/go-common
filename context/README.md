# context

This package contains all constants common to projects, including keys for working with context

For extract some value from context use:

```
import (
	"github.com/webdevelop-pro/go-common/context/keys"
    ...
)

...

requestID := keys.GetCtxValue(ctx, keys.RequestID).(string)

```