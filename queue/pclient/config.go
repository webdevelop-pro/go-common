package pclient

type Config struct {
	ServiceAccountCredentials string `split_words:"true"`
	ProjectID                 string `required:"true" split_words:"true"`
}
