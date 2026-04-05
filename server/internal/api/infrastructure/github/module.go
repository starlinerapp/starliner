package github

import (
	"github.com/palantir/go-githubapp/githubapp"
	"go.uber.org/fx"
	"starliner.app/internal/api/conf"
)

var Module = fx.Module(
	"github",
	fx.Provide(
		NewClientCreator,
		NewClient,
	),
)

func NewClientCreator(cfg *conf.Config) githubapp.ClientCreator {
	return githubapp.NewClientCreator(
		"https://api.github.com/",
		"https://api.github.com/graphql",
		cfg.GithubAppID,
		[]byte(cfg.GithubAppPrivateKey),
	)
}
