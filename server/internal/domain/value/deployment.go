package value

type DeploymentMessage struct {
	DeploymentId int64
	ClusterId    int64
	Database     Database
}
