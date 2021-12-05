package mnapi

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/go-resty/resty/v2"
)

const defaultAPIPrefix string = "/mn/api"

type RequestError struct {
	Err     error
	Status  int
	Message string `json:"error"`
}

func (re *RequestError) Error() string {
	return fmt.Sprintf("request %d error: %s", re.Status, re.Message)
}

type Client struct {
	baseURL    *url.URL
	restClient *resty.Client
}

func NewClient(mnTarget string, customHeaders map[string]string) (*Client, error) {
	client := Client{
		restClient: resty.New(),
	}
	url, err := url.Parse(mnTarget)
	if err != nil {
		return nil, err
	}
	if !url.IsAbs() {
		return nil, errors.New("mm target url must be an absolute url")
	}
	client.baseURL = url
	client.restClient.SetHostURL(mnTarget + defaultAPIPrefix)
	client.restClient.SetError(&RequestError{})
	client.restClient.SetHeaders(customHeaders)
	return &client, nil
}

func (c *Client) SetPrefix(newPrefix string) {
	c.restClient.SetHostURL(c.baseURL.String() + newPrefix)
}
