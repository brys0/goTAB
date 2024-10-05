package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Platforms struct {
	Platforms []Platform `json:"platforms"`
}

type Platform struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Type          string `json:"type"`
	Version       string `json:"version"`
	VersionID     string `json:"version_id"`
	DisplayName   string `json:"display_name"`
	ReplacementID string `json:"replacement_id"`
	Supported     bool   `json:"supported"`
	Architecture  string `json:"architecture"`
}

func (api *APIClient) FetchSupportedPlatforms() ([]Platform, error) {
	platform_url := api.Server.JoinPath("/api/v1/TestDataApi/Platforms")

	platform_request, err := http.NewRequest("GET", platform_url.String(), nil)

	if err != nil {
		return nil, err
	}

	platform_request.Header.Set("accept", "application/json")

	platform_response, err := api.Client.Do(platform_request)

	if err != nil {
		return nil, err
	}

	defer platform_response.Body.Close()

	if platform_response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid status code received: %d", platform_response.StatusCode)
	}

	platforms := &Platforms{}
	err = json.NewDecoder(platform_response.Body).Decode(&platforms)

	if err != nil {
		return nil, err
	}

	return platforms.Platforms, nil
}
