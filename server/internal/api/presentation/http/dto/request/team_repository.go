package request

type AssignRepoToTeam struct {
	GithubRepoId int64  `json:"githubRepoId" binding:"required"`
	RepoName     string `json:"repoName" binding:"required"`
}
