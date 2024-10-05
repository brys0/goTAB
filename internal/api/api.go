package api

import (
	"net/http"
	"net/url"
	"time"
)

type APIClient struct {
	Server *url.URL
	Client *http.Client
}

func CreateNewAPI(server *url.URL) *APIClient {
	return &APIClient{
		Server: server,
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}
