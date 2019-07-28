package cgc

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
)

func decodeError(r io.Reader) error {
	var resp apiErrorResponseTemplate
	if err := json.NewDecoder(r).Decode(&resp); err != nil {
		return fmt.Errorf("unmarshalling error response failed: %s", err.Error())
	}
	return fmt.Errorf("%s", resp.Message)
}

func mustParseURL(s string) *url.URL {
	if u, err := url.Parse(s); err != nil {
		panic(err)
	} else {
		return u
	}
}
