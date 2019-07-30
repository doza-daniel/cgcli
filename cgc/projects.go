package cgc

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Project struct represents the project information returned from CGC API.
type Project struct {
	Href string `json:"href"`
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Projects lists all the projects that belong to the token holder.
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
