package db

// Config is a struct to configure postgresql
type Config struct {
	DbDatabase       string `required:"true" split_words:"true"`
	DbHost           string `required:"true" split_words:"true"`
	DbPort           uint16 `default:"5432" split_words:"true"`
	DbUser           string `required:"true" split_words:"true"`
	DbPassword       string `required:"true" split_words:"true"`
	DbMaxConnections int32  `default:"16" split_words:"true"`
	DbLogLevel       string `default:"none" split_words:"true"`
	HOSTNAME         string `default:"go-application" split_words:"true"`
}
