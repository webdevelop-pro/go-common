package admin

import (
	"strconv"

	_ "github.com/GoAdminGroup/go-admin/adapter/echo"
	"github.com/GoAdminGroup/go-admin/examples/datamodel"
	_ "github.com/GoAdminGroup/go-admin/modules/db/drivers/postgres"
	"github.com/GoAdminGroup/go-admin/tests/tables"
	_ "github.com/GoAdminGroup/themes/adminlte"
	"github.com/labstack/echo/v4"
	"github.com/webdevelop-pro/go-common/db"

	"github.com/GoAdminGroup/go-admin/engine"
	"github.com/GoAdminGroup/go-admin/modules/config"
	"github.com/GoAdminGroup/go-admin/modules/language"
	"github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/chartjs"
)

type Config struct {
	Title     string `required:"true" split_words:"true"`
	Logo      string `required:"true" split_words:"true"`
	MiniLogo  string `required:"true" split_words:"true"`
	Debug     bool   `required:"true" split_words:"true"`
	Theme     string `required:"false" split_words:"true" default:"adminlte"`
	Env       string `required:"true" split_words:"true"`
	IndexUrl  string `required:"false" split_words:"true" default:"/"`
	UrlPrefix string `required:"false" split_words:"true" default:"admin"`
}

// NewAPI create api instance
func ConfigureAdmin(
	e *echo.Echo,
	conf Config,
	dbConfigs ...db.Config,
) error {
	eng := engine.Default()

	cfg := config.Config{
		Title:    conf.Title,
		Logo:     template.HTML(conf.Logo),
		MiniLogo: template.HTML(conf.MiniLogo),
		Theme:    conf.Theme,

		Env:       conf.Env,
		Databases: config.DatabaseList{},
		UrlPrefix: conf.UrlPrefix,
		IndexUrl:  conf.IndexUrl,
		Store: config.Store{
			Path:   "./uploads",
			Prefix: "uploads",
		},
		Debug:    conf.Debug,
		Language: language.CN,
	}

	for _, dbConf := range dbConfigs {
		cfg.Databases[dbConf.Database] = config.Database{
			Host:       dbConf.Host,
			Port:       strconv.Itoa(int(dbConf.Port)),
			User:       dbConf.User,
			Pwd:        dbConf.Password,
			Name:       dbConf.Database,
			MaxIdleCon: int(dbConf.DbMaxConnections),
			MaxOpenCon: int(dbConf.DbMaxConnections) * 2,
			Driver:     config.DriverPostgresql,
		}
	}

	template.AddComp(chartjs.NewChart())

	if err := eng.AddConfig(&cfg).
		AddGenerators(tables.Generators).
		AddDisplayFilterXssJsFilter().
		AddGenerator("user", datamodel.GetUserTable).
		Use(e); err != nil {
		return err
	}

	e.Static("/uploads", "./uploads")

	eng.HTML("GET", "/admin", datamodel.GetContent)

	return nil
}
