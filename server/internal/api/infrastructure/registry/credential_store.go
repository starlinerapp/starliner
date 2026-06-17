package registry

import "net/url"

type staticCredentialStore struct {
	username string
	password string
}

func (s *staticCredentialStore) Basic(*url.URL) (string, string) {
	return s.username, s.password
}

func (s *staticCredentialStore) RefreshToken(*url.URL, string) string {
	return ""
}

func (s *staticCredentialStore) SetRefreshToken(*url.URL, string, string) {}
