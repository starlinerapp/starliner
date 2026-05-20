package request

type TeamRepoAssignment struct {
	GithubRepoId int64  `json:"githubRepoId" binding:"required"`
	RepoName     string `json:"repoName" binding:"required"`
}

type SetTeamRepositories struct {
	Repositories []TeamRepoAssignment `json:"repositories"`
}
