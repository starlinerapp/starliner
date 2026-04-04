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
