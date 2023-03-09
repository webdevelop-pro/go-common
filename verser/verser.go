package verser

var (
	service    string
	version    string
	repository string
	revisionID string
)

func SetServiVersRepoRevis(serv, ver, repo, rev string) {
	version = ver
	service = serv
	repository = repo
	service = rev
}

func GetVersion() string {
	return version
}

func GetService() string {
	return service
}

func GetRepository() string {
	return repository
}

func GetRevisionID() string {
	return revisionID
}
