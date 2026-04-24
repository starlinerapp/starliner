package value

import "starliner.app/internal/api/domain/port"

type RepositoryFile struct {
	Name *string
	Path *string
	Type *string
	SHA  *string
	Size *int
	URL  string
}

func NewRepositoryFile(f *port.RepositoryFile) *RepositoryFile {
	return &RepositoryFile{
		Name: f.Name,
		Path: f.Path,
		Type: f.Type,
		SHA:  f.SHA,
		Size: f.Size,
		URL:  f.URL,
	}
}

func NewRepositoryFiles(fs []*port.RepositoryFile) []*RepositoryFile {
	files := make([]*RepositoryFile, len(fs))
	for i, f := range fs {
		files[i] = NewRepositoryFile(f)
	}
	return files
}

type PullRequestOpenedEvent struct {
	RepositoryOwner string
	RepositoryId    int64
	RepositoryName  string
	RepositoryUrl   string
	SourceBranch    string
	TargetBranch    string
	PrNumber        int
}

func (e *PullRequestOpenedEvent) EventName() string {
	return "pull_request.opened"
}

type PullRequestClosedEvent struct {
	RepositoryOwner string
	RepositoryId    int64
	RepositoryName  string
	RepositoryUrl   string
	TargetBranch    string
	PrNumber        int
	Merged          bool
}

func (e *PullRequestClosedEvent) EventName() string {
	return "pull_request.closed"
}

type PushToBranchEvent struct {
	RepositoryOwner string
	RepositoryName  string
	RepositoryUrl   string
	TargetBranch    string
	Ref             string
}

func (e *PushToBranchEvent) EventName() string {
	return "push"
}
