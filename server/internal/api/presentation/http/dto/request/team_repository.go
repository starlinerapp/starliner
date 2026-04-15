package request

type AssignRepoToTeam struct {
	GithubRepoId int64  `json:"github_repo_id" binding:"required"`
	RepoName     string `json:"repo_name" binding:"required"`
}
