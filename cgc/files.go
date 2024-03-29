package cgc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type fileStorage struct {
	Type     string `json:"type"`
	Volume   string `json:"volume"`
	Location string `json:"location"`
}

type fileOrigin struct {
	Dataset string `json:"dataset"`
}

// File struct represents the file information returned from CGC API.
type File struct {
	Project    string                 `json:"project"`
	Href       string                 `json:"href"`
	Name       string                 `json:"name"`
	ID         string                 `json:"id"`
	Size       int64                  `json:"size"`
	CreatedOn  time.Time              `json:"created_on"`
	ModifiedOn time.Time              `json:"modified_on"`
	Storage    fileStorage            `json:"storage"`
	Origin     fileOrigin             `json:"origin"`
	Tags       []string               `json:"tags"`
	Metadata   map[string]interface{} `json:"metadata"`
}

// Files lists all the files under the project with projectID.
func (c Client) Files(projectID string) ([]File, error) {
	u := mustParseURL(c.baseURL)
	u.Path += "files"
	params := url.Values{}
	params.Add("project", projectID)
	u.RawQuery = params.Encode()

	files := make([]File, 0)

	// if there's more files than returned by default by the API, links array will
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
			Files []File `json:"items"`
		}
		if err := json.NewDecoder(resp).Decode(&r); err != nil {
			return nil, fmt.Errorf("unmarshalling response failed: %s", err.Error())
		}
		files = append(files, r.Files...)

		u = nil
		for _, link := range r.Links {
			if link.Rel == "next" {
				u = mustParseURL(link.Href)
			}
		}
	}

	return files, nil
}

// StatFile gets the details of the file that has the ID of fileID.
func (c Client) StatFile(fileID string) (File, error) {
	u := mustParseURL(c.baseURL)
	u.Path += fmt.Sprintf("files/%s", fileID)
	resp, err := c.request(http.MethodGet, u, nil)
	if err != nil {
		return File{}, fmt.Errorf("fetching file details failed: %s", err.Error())
	}
	defer resp.Close()

	var file File
	if err := json.NewDecoder(resp).Decode(&file); err != nil {
		return File{}, fmt.Errorf("unmarshalling response failed: %s", err.Error())
	}

	return file, nil
}

// UpdateFile updates the file that has the ID of fileID. Updates slice represent strings with the format like
// 'key=value' or 'metadata.key=value'. Every update string gets parsed, and a request to an appropriate endpoint is made.
// If update string is used to update metadata, a PATCH request should be sent to 'files/{fileID}/metadata', else
// 'files/{fileID}'.
func (c Client) UpdateFile(fileID string, updates []string) error {
	for _, update := range updates {
		encoded, isMetadata, err := updateStringToJSON(update)
		if err != nil {
			return fmt.Errorf("encoding update string to JSON failed: %s", err.Error())
		}
		u := mustParseURL(c.baseURL)
		u.Path += fmt.Sprintf("files/%s/", fileID)
		if isMetadata {
			u.Path += "metadata/"
		}
		resp, err := c.request(http.MethodPatch, u, bytes.NewReader(encoded))
		if err != nil {
			return fmt.Errorf("updating file failed: %s", err)
		}
		defer resp.Close()
	}

	return nil
}

// DownloadFile downloads a file that has the ID of fileID and writes it to dest location on the system. Two requests have
// to be made in order to make this happen. First one get's the download URL, and the second one actually downloads the file.
func (c Client) DownloadFile(fileID, dest string) error {
	u := mustParseURL(c.baseURL)
	u.Path += fmt.Sprintf("files/%s/download_info", fileID)
	resp, err := c.request(http.MethodGet, u, nil)
	if err != nil {
		return fmt.Errorf("fetching file details failed: %s", err.Error())
	}
	defer resp.Close()

	var r struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(resp).Decode(&r); err != nil {
		return fmt.Errorf("unmarshalling response failed: %s", err.Error())
	}

	file, err := c.request(http.MethodGet, mustParseURL(r.URL), nil)
	if err != nil {
		return fmt.Errorf("download link failed: %s", err.Error())
	}
	defer file.Close()

	destf, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("creating '%s' file failed: %s", dest, err.Error())
	}

	if _, err := io.Copy(destf, file); err != nil {
		return fmt.Errorf("writing file to '%s' failed: %s", dest, err.Error())
	}

	if err := destf.Sync(); err != nil {
		return fmt.Errorf("syncing '%s' file failed: %s", dest, err.Error())
	}

	return nil
}

// updateStringToJSON parses a string in a format of 'key=value', 'metadata.key=value' and encodes it in JSON format.
// Returns true if the string is prefixed with 'metadata.'.
func updateStringToJSON(updateString string) ([]byte, bool, error) {
	kv := strings.Split(updateString, "=")
	if len(kv) != 2 {
		return nil, false, fmt.Errorf("malformed update string")
	}
	key := kv[0]
	val := kv[1]
	isMetadata := strings.HasPrefix(key, "metadata.")
	key = strings.TrimPrefix(key, "metadata.")

	// value will be interface{} so json encoding will marshal it to a correct type
	var value interface{}

	// if the key is 'tags' we expect an array of strings
	if key == "tags" {
		value = strings.Split(val, ",")
	}

	// JSON values can be boolean, numbers, arrays (we only care about arrays of strings)
	// or strings (we don't care about objects)

	// try and parse out a boolean
	b, err := strconv.ParseBool(val)
	if err == nil {
		value = b
	}
	// try and parse out a number
	f, err := strconv.ParseFloat(val, 64)
	if err == nil {
		value = f
	}
	// if all else failed, we resort to a plain string (unless the actual string is 'null')
	if val != "null" && value == nil {
		value = val
	}

	toEncode := map[string]interface{}{
		key: value,
	}
	buff := bytes.Buffer{}
	if err := json.NewEncoder(&buff).Encode(toEncode); err != nil {
		return nil, false, fmt.Errorf("encoding failed: %s", err)
	}

	return buff.Bytes(), isMetadata, nil
}
