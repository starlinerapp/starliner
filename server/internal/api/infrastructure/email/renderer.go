package email

import (
	"bytes"
	"embed"
	"html/template"
)

//go:embed template/*.html
var templatesFs embed.FS

type Renderer struct {
	tmpl *template.Template
}

func NewRenderer() *Renderer {
	tmpl, err := template.ParseFS(templatesFs, "template/*.html")
	if err != nil {
		panic(err)
	}
	return &Renderer{tmpl: tmpl}
}

func (r *Renderer) Render(name string, data any) (string, error) {
	var buf bytes.Buffer
	if err := r.tmpl.ExecuteTemplate(&buf, name, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
