package cgc

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"
)

func sampleFile() File {
	return File{
		Project:    "foo",
		Href:       "bar",
		Name:       "baz",
		ID:         "oof",
		Size:       42,
		CreatedOn:  time.Now(),
		ModifiedOn: time.Now(),
		Storage: fileStorage{
			Type:     "foo",
			Volume:   "bar",
			Location: "baz",
		},
		Origin: fileOrigin{
			Dataset: "foo",
		},
		Tags: []string{"foo", "bar", "baz"},
		Metadata: map[string]interface{}{
			"foo": 42,
			"bar": "baz",
		},
	}
}

var testFiles = []File{
	sampleFile(),
	sampleFile(),
	sampleFile(),
	sampleFile(),
}

var testProjectID = "testProjectID"
var testFileID = "testFileID"

func handleFiles(w http.ResponseWriter, r *http.Request) {
	var resp struct {
		apiOKResponseTemplate
		Items []File `json:"items"`
	}
	projectID := r.URL.Query().Get("project")
	if projectID != testProjectID {
		w.WriteHeader(http.StatusNotFound)
		resp := apiErrorResponseTemplate{
			Message: fmt.Sprintf("Project with id '%s' not found", projectID),
		}
		if err := json.NewEncoder(w).Encode(&resp); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	resp.Items = testFiles

	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func handleStatFile(w http.ResponseWriter, r *http.Request) {
	pat := regexp.MustCompile(`files/([^/]+)/?.*`)
	s := pat.FindStringSubmatch(r.URL.Path)
	if len(s) != 2 {
		w.WriteHeader(http.StatusBadRequest)
		resp := apiErrorResponseTemplate{
			Message: "Bad request",
		}
		if err := json.NewEncoder(w).Encode(&resp); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	fileID := s[1]
	if fileID != testFileID {
		w.WriteHeader(http.StatusNotFound)
		resp := apiErrorResponseTemplate{
			Message: fmt.Sprintf("File with id '%s' not found", fileID),
		}
		if err := json.NewEncoder(w).Encode(&resp); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	resp := sampleFile()
	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func TestFiles(t *testing.T) {
	type in struct {
		projectID string
	}
	type out struct {
		err error
	}

	td := []struct {
		label string
		in    in
		out   out
	}{
		{"All good", in{testProjectID}, out{nil}},
		{"Wrong Project ID", in{"wrong"}, out{errors.New("not found")}},
	}

	testToken := "test_token"
	contentType := "application/json"
	ts := httptest.NewServer(tokenMiddleware(testToken, contentTypeMiddleware(contentType, handleFiles)))
	defer ts.Close()
	client := New(testToken)
	client.baseURL = ts.URL

	for _, tt := range td {
		t.Run(tt.label, func(t *testing.T) {
			files, err := client.Files(tt.in.projectID)
			if err != nil {
				if tt.out.err != nil {
					if !strings.Contains(err.Error(), tt.out.err.Error()) {
						t.Fatalf("expected '%v', got '%v'", tt.out.err, err)
					}
					return
				}
				t.Fatalf("expected no error, got '%v'", err)
			}

			if len(files) != len(testFiles) {
				t.Fatalf("expected %d files, got %d", len(testFiles), len(files))
			}
		})
	}
}

func TestStatFile(t *testing.T) {
	type in struct {
		fileID string
	}
	type out struct {
		err error
	}
	td := []struct {
		label string
		in    in
		out   out
	}{
		{"All good", in{testFileID}, out{nil}},
		{"Wrong File ID", in{"wrong"}, out{errors.New("not found")}},
	}

	testToken := "test_token"
	contentType := "application/json"
	ts := httptest.NewServer(tokenMiddleware(testToken, contentTypeMiddleware(contentType, handleStatFile)))
	defer ts.Close()
	client := New(testToken)
	client.baseURL = ts.URL

	for _, tt := range td {
		t.Run(tt.label, func(t *testing.T) {
			_, err := client.StatFile(tt.in.fileID)
			if err != nil {
				if tt.out.err != nil {
					if !strings.Contains(err.Error(), tt.out.err.Error()) {
						t.Fatalf("expected '%v', got '%v'", tt.out.err, err)
					}
					return
				}
				t.Fatalf("expected no error, got '%v'", err)
			}
		})
	}
}

func TestUpdateStringToJSON(t *testing.T) {
	type in struct {
		updateString string
	}
	type out struct {
		body       []byte
		isMetadata bool
		err        error
	}

	td := []struct {
		label string
		in    in
		out   out
	}{
		{
			"Metadata string",
			in{"metadata.foo=bar"},
			out{
				body:       []byte(`{"foo":"bar"}`),
				isMetadata: true,
				err:        nil,
			},
		},
		{
			"Metadata int",
			in{"metadata.foo=42"},
			out{
				body:       []byte(`{"foo":42}`),
				isMetadata: true,
				err:        nil,
			},
		},
		{
			"Metadata float",
			in{"metadata.foo=42.1"},
			out{
				body:       []byte(`{"foo":42.1}`),
				isMetadata: true,
				err:        nil,
			},
		},
		{
			"Normal string",
			in{"foo=bar"},
			out{
				body:       []byte(`{"foo":"bar"}`),
				isMetadata: false,
				err:        nil,
			},
		},
		{
			"Normal int",
			in{"foo=42"},
			out{
				body:       []byte(`{"foo":42}`),
				isMetadata: false,
				err:        nil,
			},
		},
		{
			"Normal float",
			in{"foo=42.1"},
			out{
				body:       []byte(`{"foo":42.1}`),
				isMetadata: false,
				err:        nil,
			},
		},
		{
			"Malformed",
			in{"asdfj"},
			out{
				body:       nil,
				isMetadata: false,
				err:        errors.New("malformed update string"),
			},
		},
	}

	for _, tt := range td {
		t.Run(tt.label, func(t *testing.T) {
			body, isMetadata, err := updateStringToJSON(tt.in.updateString)
			if err != nil {
				if tt.out.err != nil {
					if !strings.Contains(err.Error(), tt.out.err.Error()) {
						t.Fatalf("expected '%v', got '%v'", tt.out.err, err)
					}
					return
				}
				t.Fatalf("expected no error, got '%v'", err)
			}
			if isMetadata != tt.out.isMetadata {
				t.Fatalf("expected '%s' to be metadata update string", tt.in.updateString)
			}
			bodyStr := string(body)
			expectedBodyStr := string(tt.out.body)
			if strings.TrimSpace(bodyStr) != strings.TrimSpace(expectedBodyStr) {
				t.Fatalf("expected '%s', got '%s'", expectedBodyStr, bodyStr)
			}
		})
	}
}
