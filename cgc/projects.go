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
	var r struct {
		apiOKResponseTemplate
		Items []Project `json:"items"`
	}
	if err := json.NewDecoder(resp).Decode(&r); err != nil {
		return nil, fmt.Errorf("unmarshalling response failed: %s", err.Error())
	}

	return r.Items, nil
}
