package gokick

import (
	"fmt"
	"net/http"
)

const (
	defaultAPIBaseURL = "https://api.kick.com"
	authBaseURL       = "https://id.kick.com"
)

type Client struct {
	options *ClientOptions
}

type ClientOptions struct {
	UserAccessToken string
	HTTPClient      *http.Client
	APIBaseURL      string
	AuthBaseURL     string
	ClientID        string
	ClientSecret    string
}

func NewClient(options *ClientOptions) (*Client, error) {
	if options.APIBaseURL == "" {
		options.APIBaseURL = defaultAPIBaseURL
	}

	if options.AuthBaseURL == "" {
		options.AuthBaseURL = authBaseURL
	}

	if options.HTTPClient == nil {
		options.HTTPClient = &http.Client{}
	}

	return &Client{options: options}, nil
}

type errorResponse struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func (c *Client) buildURL(base, path string) string {
	return fmt.Sprintf("%s%s", base, path)
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.options.UserAccessToken))

	response, err := c.options.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	return response, nil
}
