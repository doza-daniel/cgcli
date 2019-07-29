package cgc

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func tokenMiddleware(token string, f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(tokenHeader) != token {
			w.WriteHeader(http.StatusUnauthorized)
			x := apiErrorResponseTemplate{
				Message: "Unauthorized",
			}
			if err := json.NewEncoder(w).Encode(&x); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		} else {
			f(w, r)
		}
	}
}

func contentTypeMiddleware(contentType string, f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bs, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err := r.Body.Close(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if r.Method != http.MethodGet && (r.Header.Get("Content-Type") != contentType || !json.Valid(bs)) {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			x := apiErrorResponseTemplate{
				Message: "Unsupported Media Type.",
			}
			if err := json.NewEncoder(w).Encode(&x); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			r.Body = ioutil.NopCloser(bytes.NewBuffer(bs))
			f(w, r)
		}
	}
}

func TestRequest(t *testing.T) {
	testToken := "asdf"
	type in struct {
		method string
		token  string
		body   io.Reader
	}
	type out struct {
		data []byte
		err  error
	}
	td := []struct {
		label string
		in    in
		out   out
	}{
		{
			"All good",
			in{http.MethodGet, testToken, nil},
			out{[]byte("OK"), nil},
		},
		{
			"Unauthorized",
			in{http.MethodGet, "wrong_token", nil},
			out{[]byte(""), errors.New("Unauthorized")},
		},
		{
			"Bad Json",
			in{http.MethodPatch, testToken, strings.NewReader("not json")},
			out{[]byte(""), errors.New("Unsupported Media Type")},
		},
		{
			"Good Json",
			in{http.MethodPatch, testToken, strings.NewReader(`{"valid":"json"}`)},
			out{[]byte("OK"), nil},
		},
	}

	f := func(w http.ResponseWriter, r *http.Request) {
		io.Copy(w, strings.NewReader("OK"))
	}
	ts := httptest.NewServer(tokenMiddleware(testToken, contentTypeMiddleware("application/json", f)))

	for _, tt := range td {
		t.Run(tt.label, func(t *testing.T) {
			c := New(tt.in.token)
			c.baseURL = ts.URL

			resp, err := c.request(tt.in.method, mustParseURL(ts.URL), tt.in.body)
			if err != nil {
				if tt.out.err != nil {
					if !strings.Contains(err.Error(), tt.out.err.Error()) {
						t.Fatalf("got '%v', want '%v'", err.Error(), tt.out.err.Error())
					}
					return
				}
				t.Fatalf("expected no error, got '%s'", err.Error())
			}
			defer resp.Close()
			output, err := ioutil.ReadAll(resp)
			if err != nil {
				t.Fatalf("reading bytes from output failed: %s", err.Error())
			}
			if bytes.Compare(tt.out.data, output) != 0 {
				t.Fatalf("got '%s', wanted '%s'", string(output), string(tt.out.data))
			}
		})
	}
}
