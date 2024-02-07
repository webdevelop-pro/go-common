package broker

type Config struct {
	ServiceAccountCredentials string `required:"true" split_words:"true"`
	ProjectID                 string `required:"true" split_words:"true"`
	Topic                     string `required:"true"`
	Subscription              string `required:"true"`
}
