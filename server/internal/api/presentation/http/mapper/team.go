package mapper

import (
	"starliner.app/internal/api/domain/value"
	"starliner.app/internal/api/presentation/http/dto/request"
)

func MapTeamReposFromRequest(
	teamID int64,
	repos []request.TeamRepoAssignment,
) []*value.TeamRepo {
	result := make([]*value.TeamRepo, len(repos))
	for i, repo := range repos {
		result[i] = &value.TeamRepo{
			TeamId:       teamID,
			GithubRepoId: repo.GithubRepoId,
			RepoName:     repo.RepoName,
		}
	}
	return result
}
