package esign

type config struct {
	host          string
	projectId     string
	projectSecret string
}

var conf = new(config)

func Config() *config {
	return conf
}

func (e *config) Host() string {
	return e.host
}

func (e *config) SetHost(host string) {
	e.host = host
}

func (e *config) ProjectId() string {
	return e.projectId
}

func (e *config) SetProjectId(projectId string) {
	e.projectId = projectId
}

func (e *config) ProjectSecret() string {
	return e.projectSecret
}

func (e *config) SetProjectSecret(projectSecret string) {
	e.projectSecret = projectSecret
}
