package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/brys0/goTAB/internal/hardware"
)

type TestRoot struct {
	Token  string              `json:"token"`
	FFmpeg hardware.FFmpegInfo `json:"ffmpeg"`
	Tests  []Test              `json:"tests"`
}

type Test struct {
	Name         string                `json:"name"`
	SourceURL    string                `json:"source_url"`
	SourceHashes []hardware.HashedFile `json:"source_hashs"`
	TestType     string                `json:"test_type"`
	Data         []TestData            `json:"data"`
}

type TestData struct {
	ID             string             `json:"id"`
	FromResolution string             `json:"from_resolution"`
	ToResolution   string             `json:"to_resolution"`
	Bitrate        int                `json:"bitrate"`
	Arguments      []TestDataArgument `json:"arguments"`
}

type TestDataArgument struct {
	Type      string `json:"type"`
	Arguments string `json:"args"`
	Codec     string `json:"codec"`
}

func (api *APIClient) FetchPlatformTests(os_id string) (*TestRoot, error) {
	test_url := api.Server.JoinPath("/api/v1/TestDataApi")

	test_url_query := test_url.Query()
	test_url_query.Add("platformId", os_id)
	test_url.RawQuery = test_url_query.Encode()

	test_request, err := http.NewRequest("GET", test_url.String(), nil)

	if err != nil {
		return nil, err
	}

	test_request.Header.Set("accept", "application/json")

	test_response, err := api.Client.Do(test_request)

	if err != nil {
		return nil, err
	}

	defer test_response.Body.Close()

	if test_response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid status code received: %d", test_response.StatusCode)
	}

	test := &TestRoot{}
	err = json.NewDecoder(test_response.Body).Decode(test)

	if err != nil {
		return nil, err
	}

	return test, nil
}
