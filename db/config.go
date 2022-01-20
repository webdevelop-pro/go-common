package db

// Type represent storage engine type
type Type string
type Active string

var (
	Postgres Type   = "postgres"
	MySQL    Type   = "mysql"
	SqlLite  Type   = "sqllite"
	Enable   Active = "enable"
	Disable  Active = "disable"
)

// Config is a struct to configure postgresql
type Config struct {
	Type           Type   `required:"true" split_words:"true"`
	Host           string `required:"false" default:"localhost" split_words:"true"`
	Port           uint16 `default:"5432" split_words:"true"`
	User           string `required:"true" split_words:"true"`
	Password       string `required:"true" split_words:"true"`
	Database       string `required:"true" split_words:"true"`
	SSLMode        Active `required:"true" split_words:"true"`
	MaxConnections int    `default:"16" split_words:"true"`
	HOSTNAME       string `default:"go-application" split_words:"true"`

	// TODO: Why we use DB prefix? I suggest to remove db prefix from all fileds because it's config for db
	// You can remove code below if your are agree, and fix
	DbDatabase       string `required:"true" split_words:"true"`
	DbHost           string `required:"true" split_words:"true"`
	DbPort           uint16 `default:"5432" split_words:"true"`
	DbUser           string `required:"true" split_words:"true"`
	DbPassword       string `required:"true" split_words:"true"`
	DbMaxConnections int32  `default:"16" split_words:"true"`
	DbLogLevel       string `default:"none" split_words:"true"`
}
