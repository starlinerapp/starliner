package ansible

type Output struct {
	Plays []struct {
		Tasks []struct {
			Hosts map[string]struct {
				Content string `json:"content,omitempty"`
			} `json:"hosts"`
		} `json:"tasks"`
	} `json:"plays"`
}
