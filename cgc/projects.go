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

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("creating projects request failed: %s", err.Error())
	}
	req.Header.Add(tokenHeader, c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("projects request failed: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code: %d, message: %s", resp.StatusCode, decodeError(resp.Body))
	}

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
	if err := json.NewDecoder(resp.Body).Decode(&respJSON); err != nil {
		return nil, fmt.Errorf("unmarshalling response failed: %s", err.Error())
	}

	return respJSON.Items, nil
}
