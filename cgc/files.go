package cgc

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

// File ...
type File struct {
	Project    string    `json:"project"`
	Href       string    `json:"href"`
	Name       string    `json:"name"`
	ID         string    `json:"id"`
	Size       int64     `json:"size"`
	CreatedOn  time.Time `json:"created_on"`
	ModifiedOn time.Time `json:"modified_on"`
	Storage    struct {
		Type     string `json:"type"`
		Volume   string `json:"volume"`
		Location string `json:"location"`
	} `json:"storage"`
	Origin struct {
		Dataset string `json:"dataset"`
	} `json:"origin"`
	Tags     []string `json:"tags"`
	Metadata metadata `json:"metadata"`
}

type metadata struct {
	metadataCategoryFile
	metadataCategorySample
	metadataCategoryAliquot
	metadataCategoryCase
	metadataCategoryGeneral
}

type metadataCategoryFile struct {
	LibraryID            string `json:"library_id"`
	Platform             string `json:"platform"`
	PlatformUnitID       string `json:"platform_unit_id"`
	PairedEnd            string `json:"paired_end"`
	FileSegmentNumber    int    `json:"file_segment_number"`
	QualityScale         string `json:"quality_scale"`
	ExperimentalStrategy string `json:"experimental_strategy"`
	ReferenceGenome      string `json:"reference_genome"`
}

type metadataCategorySample struct {
	SampleID   string `json:"sample_id"`
	SampleType string `json:"sample_type"`
	SampleUUID string `json:"sample_uuid"`
}

type metadataCategoryAliquot struct {
	AliquotID   string `json:"aliquot_id"`
	AliquotUUID string `json:"aliquot_uuid"`
}

type metadataCategoryCase struct {
	CaseID         string `json:"case_id"`
	CaseUUID       string `json:"case_uuid"`
	PrimarySite    string `json:"primary_site"`
	DiseaseType    string `json:"disease_type"`
	Gender         string `json:"gender"`
	AgeAtDiagnosis uint   `json:"age_at_diagnosis"`
	VitalStatus    string `json:"vital_status"`
	DaysToDeath    uint   `json:"days_to_death"`
	Race           string `json:"race"`
	Ethnicity      string `json:"ethnicity"`
}

type metadataCategoryGeneral struct {
	Investigation string `json:"investigation"`
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

// StatFile ...
func (c Client) StatFile(fileID string) (File, error) {
	u := mustParseURL(c.baseURL)
	u.Path += fmt.Sprintf("files/%s", fileID)
	resp, err := c.get(u)
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
