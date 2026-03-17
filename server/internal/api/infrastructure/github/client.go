package github

import (
	"context"
	"github.com/google/go-github/v84/github"
	"github.com/palantir/go-githubapp/githubapp"
	"starliner.app/internal/api/domain/port"
)

type Client struct {
	creator githubapp.ClientCreator
}

func NewClient(creator githubapp.ClientCreator) *Client {
	return &Client{
		creator: creator,
	}
}

func (c *Client) installationClient(installationId int64) (*github.Client, error) {
	gh, err := c.creator.NewInstallationClient(installationId)
	if err != nil {
		return nil, err
	}

	return gh, nil
}

func (c *Client) ListRepositories(ctx context.Context, installationId int64) ([]*port.Repository, error) {
	gh, err := c.installationClient(installationId)
	if err != nil {
		return nil, err
	}

	var all []*port.Repository
	opts := &github.ListOptions{
		PerPage: 100,
	}

	for {
		repos, resp, err := gh.Apps.ListRepos(ctx, opts)
		if err != nil {
			return nil, err
		}

		for _, r := range repos.Repositories {
			all = append(all, &port.Repository{
				Id:          r.ID,
				Name:        r.Name,
				FullName:    r.FullName,
				Description: r.Description,
				CreatedAt:   r.CreatedAt.GetTime(),
				PushedAt:    r.PushedAt.GetTime(),
				UpdatedAt:   r.UpdatedAt.GetTime(),
				CloneURL:    r.CloneURL,
			})
		}

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return all, nil
}
