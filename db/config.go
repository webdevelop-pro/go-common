package db

// Type represent storage engine type
type (
	Type   string
	Active string
)

var (
	Postgres Type   = "postgres"
	MySQL    Type   = "mysql"
	SQLLite  Type   = "sqllite"
	Enable   Active = "enable"
	Disable  Active = "disable"
)

// Config is a struct to configure postgresql
type Config struct {
	Type            Type   `required:"true" split_words:"true"`
	Host            string `required:"false" default:"localhost" split_words:"true"`
	Port            uint16 `default:"5432" split_words:"true"`
	User            string `required:"true" split_words:"true"`
	Password        string `required:"true" split_words:"true"`
	Database        string `required:"true" split_words:"true"`
	SSLMode         Active `required:"true" split_words:"true"`
	AppName         string `required:"true" split_words:"true"`
	MaxConnections  int    `default:"16" split_words:"true"`
	MaxConnLifetime int    `default:"3600" split_words:"true"`
	LogLevel        string `default:"error" split_words:"true"`
}
