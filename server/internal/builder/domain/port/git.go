package port

type Git interface {
	CloneRepository(repoUrl string) (dir string, commitHash string, err error)
}
