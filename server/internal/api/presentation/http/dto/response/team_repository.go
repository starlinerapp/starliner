package response

import "starliner.app/internal/api/domain/value"

type TeamRepo struct {
	TeamId       int64  `json:"teamId" binding:"required"`
	GithubRepoId int64  `json:"githubRepoId" binding:"required"`
	RepoName     string `json:"repoName" binding:"required"`
}

func NewTeamRepo(tr *value.TeamRepo) TeamRepo {
	return TeamRepo{
		TeamId:       tr.TeamId,
		GithubRepoId: tr.GithubRepoId,
		RepoName:     tr.RepoName,
	}
}
func NewTeamRepos(trs []*value.TeamRepo) []TeamRepo {
	res := make([]TeamRepo, len(trs))
	for i, tr := range trs {
		res[i] = NewTeamRepo(tr)
	}
	return res
}
