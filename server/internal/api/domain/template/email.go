package template

import (
	"bytes"
	"html/template"
	"starliner.app/internal/api/domain/template/files"
)

var inviteTmpl = template.Must(template.New("invite").Parse(files.InviteEmail))

type InviteData struct {
	OrganizationName string
	InviteLink       string
}

func RenderInvite(data InviteData) (string, error) {
	var buf bytes.Buffer
	if err := inviteTmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
