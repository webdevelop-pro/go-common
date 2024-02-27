package queue

type Config struct {
	ServiceAccountCredentials string `required:"true" split_words:"true"`
	ProjectID                 string `required:"true" split_words:"true"`
}
