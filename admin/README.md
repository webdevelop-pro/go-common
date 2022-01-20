
# GoAdmin Instruction

GoAdmin is a golang framework help gopher quickly build a data visualization platform.

- [github](https://github.com/GoAdminGroup/go-admin)
- [forum](http://discuss.go-admin.com)
- [document](https://book.go-admin.cn)

## Integrate with echo framework

If you have service based on echo framework you can easily integrate admin panel to it

Example:
```
import (
	"github.com/your-package/tables"
    "github.com/webdevelop-pro/go-common/admin"
    "github.com/webdevelop-pro/go-common/db"
)

type config struct {
	DB db.Config `required:"true" split_words:"true"`

	AdminPanel admin.Config `required:"true" split_words:"true"`
}

func main() {
    var conf config
    e := echo.New()

    ...

    err := admin.ConfigureAdmin(e, conf.AdminPanel, tables.Generators, conf.DB)
    if err != nil {
        panic(err)
    }

    ...

    e.Start(...)
}
```

## use adm cmd tool

Install via:

```
go install github.com/GoAdminGroup/adm
```

Generate/Update tables code from existing database

```
adm generate
```
