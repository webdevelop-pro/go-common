package server

// Config is struct to configure HTTP server
type Config struct {
	Host  string `required:"true"`
	Port  string `required:"true"`
	Debug bool   `default:"false"`
}
