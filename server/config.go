package server

// Config is struct to configure HTTP server
type Config struct {
	Host  string `default:""`
	Port  string `default:"8085"`
	Debug bool   `default:"false"`
}
