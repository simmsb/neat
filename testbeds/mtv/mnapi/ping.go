package mnapi

import (
	"fmt"
	"strings"
)

type PingData struct {
	Sender   string  `json:"sender"`
	Target   string  `json:"target"`
	Sent     int     `json:"sent"`
	Received int     `json:"received"`
	AvgRTT   float64 `json:"rtt_avg"`
}

func (c *Client) PingAll() ([]*PingData, error) {
	var pingData []*PingData
	resp, err := c.restClient.R().
		SetHeader("Accept", "application/json").
		SetResult(&pingData).Get("/pingall")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("received non-200 status code (%d)", resp.StatusCode())
	}
	return pingData, nil
}

func (c *Client) PingSet(nodes []string) (map[string]*PingData, error) {
	nodeParam := strings.Join(nodes, ",")
	var pingData map[string]*PingData
	resp, err := c.restClient.R().
		SetHeader("Accept", "application/json").
		SetResult(&pingData).SetQueryParam("hosts", nodeParam).
		Get("/pingset")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("received non-200 status code (%d)", resp.StatusCode())
	}
	return pingData, nil

}

func (c *Client) PingFrom(host string, dst string) (map[string]*PingData, error) {
	var pingData map[string]*PingData
	resp, err := c.restClient.R().
		SetHeader("Accept", "application/json").
		SetResult(&pingData).SetQueryParam("host", host).SetQueryParam("dst", dst).
		Get("/pingfrom")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("received non-200 status code (%d)", resp.StatusCode())
	}
	return pingData, nil

}
