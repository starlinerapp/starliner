package port

type Git interface {
	CloneRepository(repoUrl string, branchName string, accessToken string) (dir string, commitHash string, err error)
}
