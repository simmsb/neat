package mnapi

import "fmt"

type NodeInfo struct {
	Name  string   `json:"name"`
	Class string   `json:"class"`
	IPs   []string `json:"ips,omitempty"`
	MACs  []string `json:"macs,omitempty"`
}

func (c *Client) GetNodes() (map[string][]string, error) {
	// var nodes map[string]*NodeInfo
	var nodes map[string][]string
	resp, err := c.restClient.R().
		SetHeader("Accept", "application/json").
		SetResult(&nodes).Get("/nodes")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("received non-200 status code (%d)", resp.StatusCode())
	}
	return nodes, nil
}

func (c *Client) GetNodesOfClass(class string) (map[string]*NodeInfo, error) {
	var nodes map[string]*NodeInfo
	resp, err := c.restClient.R().
		SetHeader("Accept", "application/json").
		SetResult(&nodes).SetQueryParam("class", class).
		Get("/nodes")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("received non-200 status code (%d)", resp.StatusCode())
	}
	return nodes, nil
}

func (c *Client) GetNodeInfo(nodeName string) (*NodeInfo, error) {
	var node *NodeInfo
	resp, err := c.restClient.R().
		SetHeader("Accept", "application/json").
		SetResult(&node).SetPathParam("node_name", nodeName).
		Get("/node/{node_name}")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("received non-200 status code (%d)", resp.StatusCode())
	}
	return node, nil
}
