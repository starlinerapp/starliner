package git

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"os"
	"starliner.app/internal/builder/domain/port"
)

type Git struct {
}

func NewGit() port.Git {
	return &Git{}
}

func (g *Git) CloneRepository(repoUrl string) (dir string, commitHash string, err error) {
	dir, err = os.MkdirTemp("", "repo-*")
	if err != nil {
		return "", "", err
	}

	repo, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL:           repoUrl,
		ReferenceName: plumbing.NewBranchReferenceName("main"),
		SingleBranch:  true,
		Depth:         1, // doesn't download the full commit history
	})
	if err != nil {
		err := os.RemoveAll(dir)
		if err != nil {
			return "", "", err
		}

		return "", "", err
	}
	ref, err := repo.Head()
	if err != nil {
		_ = os.RemoveAll(dir)
		return "", "", err
	}

	hash := ref.Hash().String()
	return dir, hash, nil
}
