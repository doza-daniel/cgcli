package cgc

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Project ...
type Project struct {
	Href string `json:"href"`
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Projects ...
func (c Client) Projects() ([]Project, error) {
	u := mustParseURL(c.baseURL)
	u.Path += "projects"

	resp, err := c.request(http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("fetching files failed: %s", err.Error())
	}
	defer resp.Close()

	// TODO: paging
	var respJSON struct {
		Href  string    `json:"href"`
		Items []Project `json:"items"`
		Links []struct {
			Href   string `json:"href"`
			Rel    string `json:"rel"`
			Method string `json:"method"`
		} `json:"links"`
	}
	if err := json.NewDecoder(resp).Decode(&respJSON); err != nil {
		return nil, fmt.Errorf("unmarshalling response failed: %s", err.Error())
	}

	return respJSON.Items, nil
}
