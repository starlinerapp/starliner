package port

type Git interface {
	CloneRepository(repoUrl string, accessToken string) (dir string, commitHash string, err error)
}
