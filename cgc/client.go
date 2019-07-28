package cgc

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	tokenHeader = "X-SBG-Auth-Token"
	baseURL     = "https://cgc-api.sbgenomics.com/v2/"
)

type apiErrorResponseTemplate struct {
	Message string `json:"message"`
}

type apiOKResponseTemplate struct {
	Href  string `json:"href"`
	Links []struct {
		Href   string `json:"href"`
		Rel    string `json:"rel"`
		Method string `json:"method"`
	} `json:"links"`
}

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

func (c Client) request(method string, u *url.URL, body io.Reader) (io.ReadCloser, error) {
	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, fmt.Errorf("creating request failed: %s", err.Error())
	}
	req.Header.Add(tokenHeader, c.token)
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %s", err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %d, message: %s", resp.StatusCode, decodeError(resp.Body))
	}

	return resp.Body, nil
}
