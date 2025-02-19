package gokick

import (
	"fmt"
	"net/http"
)

const defaultAPIBaseURL = "https://api.kick.com"

type Client struct {
	innerClient *http.Client
	url         string
	accessToken string
}

func NewClient(client *http.Client, url, accessToken string) (*Client, error) {
	if url == "" {
		url = defaultAPIBaseURL
	}

	return &Client{innerClient: client, url: url, accessToken: accessToken}, nil
}

type errorResponse struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func (c *Client) buildURL(path string) string {
	return fmt.Sprintf("%s%s", c.url, path)
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.accessToken))

	response, err := c.innerClient.Do(req)
	if err != nil {
		return nil, err
	}

	return response, nil
}
