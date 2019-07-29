package cgc

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

var testItems = []Project{
	Project{"asdf", "fdsa", "zcvb"},
	Project{"asdf", "fdsa", "zcvb"},
	Project{"asdf", "fdsa", "zcvb"},
	Project{"asdf", "fdsa", "zcvb"},
}

func handleProjects(w http.ResponseWriter, r *http.Request) {
	var resp struct {
		apiOKResponseTemplate
		Items []Project `json:"items"`
	}

	resp.Items = testItems

	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func TestProjects(t *testing.T) {
	testToken := "test_token"
	contentType := "application/json"
	ts := httptest.NewServer(tokenMiddleware(testToken, contentTypeMiddleware(contentType, handleProjects)))
	defer ts.Close()

	client := New(testToken)
	client.baseURL = ts.URL

	projects, err := client.Projects()
	if err != nil {
		t.Fatalf("fetching projects errored: %s", err.Error())
	}
	if len(projects) != len(testItems) {
		t.Fatalf("expected %d projects, got %d", len(testItems), len(projects))
	}
}
