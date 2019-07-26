package cgc

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// File ...
type File struct {
	Project string `json:"project"`
	Href    string `json:"href"`
	Name    string `json:"name"`
	ID      string `json:"id"`
}

// Files ...
func (c Client) Files(projectID string) ([]File, error) {
	u := mustParseURL(c.baseURL)
	u.Path += "files"
	params := url.Values{}
	params.Add("project", projectID)
	u.RawQuery = params.Encode()

	resp, err := c.get(u)
	if err != nil {
		return nil, fmt.Errorf("fetching files failed: %s", err.Error())
	}
	defer resp.Close()

	// TODO: paging
	var respJSON struct {
		Href  string `json:"href"`
		Items []File `json:"items"`
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
