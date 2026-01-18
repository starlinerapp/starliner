package ansible

import _ "embed"

//go:embed playbook/k3s.yaml
var K3sPlaybook string
