package entity

type Environment struct {
	Id        int64
	Slug      string
	Name      string
	Namespace string
}

type PreviewEnvironment struct {
	Id                 int64
	Slug               string
	Name               string
	Namespace          string
	GithubRepositoryId int64
	PrNumber           int64
}
