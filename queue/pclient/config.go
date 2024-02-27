package pclient

type Config struct {
	ServiceAccountCredentials string `required:"true" split_words:"true"`
	ProjectID                 string `required:"true" split_words:"true"`
	Topic                     string ``
	Subscription              string ``
}
