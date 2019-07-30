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

	projects := make([]Project, 0)

	// if there's more projects than returned by default by the API, links array will
	// be provided. The object that has the 'rel' field with the value of 'next' will
	// also contain the 'href' with the complete link to the next page.
	for u != nil {
		resp, err := c.request(http.MethodGet, u, nil)
		if err != nil {
			return nil, fmt.Errorf("fetching files failed: %s", err.Error())
		}
		defer resp.Close()

		var r struct {
			apiOKResponseTemplate
			Projects []Project `json:"items"`
		}
		if err := json.NewDecoder(resp).Decode(&r); err != nil {
			return nil, fmt.Errorf("unmarshalling response failed: %s", err.Error())
		}
		projects = append(projects, r.Projects...)

		u = nil
		for _, link := range r.Links {
			if link.Rel == "next" {
				u = mustParseURL(link.Href)
			}
		}
	}

	return projects, nil
}
