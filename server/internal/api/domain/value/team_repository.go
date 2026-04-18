package value

import "starliner.app/internal/api/domain/entity"

type TeamRepo struct {
	TeamId       int64
	GithubRepoId int64
	RepoName     string
}

func NewTeamRepo(tr *entity.TeamRepository) *TeamRepo {
	return &TeamRepo{
		TeamId:       tr.TeamId,
		GithubRepoId: tr.GithubRepoId,
		RepoName:     tr.RepoName,
	}
}
func NewTeamRepos(trs []*entity.TeamRepository) []*TeamRepo {
	repos := make([]*TeamRepo, len(trs))
	for i, tr := range trs {
		repos[i] = NewTeamRepo(tr)
	}
	return repos
}
