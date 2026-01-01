package config

func (s *Service) GetAcmeUrl() (string, string) {
	acmeUrl := s.config.ACME.ACMEUrl
	email := s.config.ACME.Email
	if s.config.ACME.UseStaging {
		acmeUrl = "https://acme-staging-v02.api.letsencrypt.org/directory"
	} else if s.config.ACME.ACMEUrl == "" {
		acmeUrl = "https://acme-v02.api.letsencrypt.org/directory"
	}

	if s.config.ACME.Email == "" {
		email = "cloudness@localhost.com"
	}

	return acmeUrl, email
}
