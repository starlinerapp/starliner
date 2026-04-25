package response

import "starliner.app/internal/api/domain/value"

type RepositoryFile struct {
	Name *string `json:"name" binding:"required"`
	Path *string `json:"path" binding:"required"`
	Type *string `json:"type" binding:"required"`
	SHA  *string `json:"sha" binding:"required"`
	Size *int    `json:"size" binding:"required"`
	URL  string  `json:"url" binding:"required"`
}

func NewRepositoryFile(file *value.RepositoryFile) RepositoryFile {
	return RepositoryFile{
		Name: file.Name,
		Path: file.Path,
		Type: file.Type,
		SHA:  file.SHA,
		Size: file.Size,
		URL:  file.URL,
	}
}

func NewRepositoryFiles(files []*value.RepositoryFile) []RepositoryFile {
	res := make([]RepositoryFile, len(files))
	for i, f := range files {
		res[i] = NewRepositoryFile(f)
	}
	return res
}

type FileContent struct {
	Content string `json:"content"`
}

func NewFileContent(content string) FileContent {
	return FileContent{Content: content}
}
