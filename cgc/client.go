package cgc

import (
	"net/http"
)

const (
	tokenHeader = "X-SBG-Auth-Token"
	baseURL     = "https://cgc-api.sbgenomics.com/v2/"
)

// Client ...
type Client struct {
	token      string
	httpClient *http.Client
	baseURL    string
}

// New returns an initialized CGC client.
func New(token string) Client {
	return Client{
		token:      token,
		httpClient: http.DefaultClient,
		baseURL:    baseURL,
	}
}
