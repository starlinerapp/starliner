package port

type Git interface {
	CloneRepository(repoUrl string) (string, error)
}
