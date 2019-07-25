package cgc

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
)

func decodeError(r io.Reader) error {
	var respJSON struct {
		Message string `json:"message"`
	}
	if err := json.NewDecoder(r).Decode(&respJSON); err != nil {
		return fmt.Errorf("unmarshalling error response failed: %s", err.Error())
	}
	return fmt.Errorf("%s", respJSON.Message)
}

func mustParseURL(s string) *url.URL {
	if u, err := url.Parse(s); err != nil {
		panic(err)
	} else {
		return u
	}
}
