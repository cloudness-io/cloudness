package types

import "net/url"

type Netrc struct {
	Machine  string `json:"machine"`
	Login    string `json:"login"`
	Password string `json:"password"`
	HttpUrl  string `json:"http_url"`
}

// SetMachine sets the netrc machine from a URL value.
func (n *Netrc) SetMachine(address string) error {
	url, err := url.Parse(address)
	if err != nil {
		return err
	}
	n.Machine = url.Hostname()
	n.HttpUrl = address
	return nil
}
