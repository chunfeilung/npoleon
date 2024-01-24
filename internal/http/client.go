package http

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

// ----------------------------------------------------------------------------

type ClientInterface interface {
	Fetch(url string) ([]byte, error)
}

// ----------------------------------------------------------------------------

type Client struct {
}

func (c Client) Fetch(url string) ([]byte, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "npoleon")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// ----------------------------------------------------------------------------

type FakeClient struct {
	Responses map[string][]byte
}

func (fc FakeClient) MakeFetchReturn(url string, response string) {
	fc.Responses[url] = []byte(response)
}

func (fc FakeClient) Fetch(url string) ([]byte, error) {
	if _, exists := fc.Responses[url]; exists {
		return fc.Responses[url], nil
	}

	msg := fmt.Sprintf("No response defined for endpoint '%s'", url)
	return nil, errors.New(msg)
}
