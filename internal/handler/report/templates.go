package report

type ProjectTemplate struct {
	ID             string
	LastUpdateTime string
}

type IndexTemplate struct {
	Projects []ProjectTemplate
}
