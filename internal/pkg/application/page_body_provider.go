package application

import (
	"io"
	"net/http"
)

func NewResourceBodyProvider(client *http.Client) *ResourceBodyProvider {
	return &ResourceBodyProvider{client: client}
}

type ResourceBodyProvider struct {
	client *http.Client
}

func (pbp *ResourceBodyProvider) GetResourceBody(url string) (io.ReadCloser, error) {
	client := http.Client{}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}
