package port

type DatabaseCredentials struct {
	DatabaseName string
	Username     string
	Password     string
}
type Secret interface {
	GetDatabaseCredentials(releaseName string, kubeconfigBase64 string) (*DatabaseCredentials, error)
}
