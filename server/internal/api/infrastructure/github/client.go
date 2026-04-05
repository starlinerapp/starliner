package github

import (
	"context"
	"fmt"
	"github.com/google/go-github/v84/github"
	"github.com/palantir/go-githubapp/githubapp"
	"starliner.app/internal/api/domain/port"
	"starliner.app/internal/api/domain/value"
	"strings"
)

type Client struct {
	creator githubapp.ClientCreator
}

func NewClient(creator githubapp.ClientCreator) port.GitHub {
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

func (c *Client) appClient() (*github.Client, error) {
	gh, err := c.creator.NewAppClient()
	if err != nil {
		return nil, err
	}

	return gh, nil
}

func (c *Client) GetInstallationToken(ctx context.Context, installationId int64) (string, error) {
	gh, err := c.appClient()
	if err != nil {
		return "", err
	}

	token, _, err := gh.Apps.CreateInstallationToken(ctx, installationId, nil)
	if err != nil {
		return "", err
	}

	if token.Token == nil {
		return "", fmt.Errorf("failed to create installation token")
	}
	return *token.Token, nil
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
				Owner:       r.Owner.Login,
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

func (c *Client) ListRepositoryContents(
	ctx context.Context,
	installationId int64,
	owner string,
	repository string,
	path string,
) ([]*port.RepositoryFile, error) {
	gh, err := c.installationClient(installationId)
	if err != nil {
		return nil, err
	}

	_, dirContents, _, err := gh.Repositories.GetContents(ctx, owner, repository, path, nil)
	if err != nil {
		return nil, err
	}

	var result []*port.RepositoryFile

	for _, item := range dirContents {
		result = append(result, &port.RepositoryFile{
			Name: item.Name,
			Path: item.Path,
			Type: item.Type,
			SHA:  item.SHA,
			Size: item.Size,
			URL:  item.GetHTMLURL(),
		})
	}

	return result, nil
}

func (c *Client) ParseGitEvent(eventType string, eventPayload []byte) (port.GitEvent, error) {
	switch eventType {
	case "pull_request":
		event, err := github.ParseWebHook(eventType, eventPayload)
		if err != nil {
			return nil, err
		}

		prEvent, ok := event.(*github.PullRequestEvent)
		if !ok {
			return nil, fmt.Errorf("unexpected event type: %T", event)
		}
		if prEvent.GetAction() != "closed" {
			return nil, nil
		}

		return &value.PullRequestClosedEvent{
			RepositoryOwner: prEvent.GetRepo().GetOwner().GetLogin(),
			RepositoryName:  prEvent.GetRepo().GetName(),
			RepositoryUrl:   prEvent.GetRepo().GetCloneURL(),
			TargetBranch:    prEvent.GetPullRequest().GetBase().GetRef(),
			Merged:          prEvent.GetPullRequest().GetMerged(),
		}, nil

	case "push":
		event, err := github.ParseWebHook(eventType, eventPayload)
		if err != nil {
			return nil, err
		}

		pushEvent, ok := event.(*github.PushEvent)
		if !ok {
			return nil, fmt.Errorf("unexpected event type: %T", event)
		}

		return &value.PushToBranchEvent{
			RepositoryOwner: pushEvent.GetRepo().GetOwner().GetName(),
			RepositoryName:  pushEvent.GetRepo().GetName(),
			RepositoryUrl:   pushEvent.GetRepo().GetCloneURL(),
			TargetBranch:    strings.TrimPrefix(pushEvent.GetRef(), "refs/heads/"),
			Ref:             pushEvent.GetRef(),
		}, nil

	default:
		return nil, fmt.Errorf("unsupported event type: %s", eventType)
	}
}
