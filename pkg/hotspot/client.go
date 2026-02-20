package hotspot

import (
	"fmt"

	routeros "github.com/go-routeros/routeros/v3"
)

// Client interface for MikroTik Hotspot Management
type Client struct {
	routerID uint
	client   *routeros.Client
	config   *Config
}

// Config holds configuration for hotspot client
type Config struct {
	RouterID    uint
	HotspotName string
	Currency    string
	Debug       bool
}

// NewClient creates new hotspot client
func NewClient(routerID uint, mikrotikClient *routeros.Client) *Client {
	return &Client{
		routerID: routerID,
		client:   mikrotikClient,
	}
}

// NewClientWithConfig creates client with custom config
func NewClientWithConfig(mikrotikClient *routeros.Client, config *Config) *Client {
	return &Client{
		routerID: config.RouterID,
		client:   mikrotikClient,
		config:   config,
	}
}

// GetRouterID returns the router ID
func (c *Client) GetRouterID() uint {
	return c.routerID
}

// Execute command with error handling
func (c *Client) execute(cmd string, args ...string) (*routeros.Reply, error) {
	fullArgs := append([]string{cmd}, args...)
	reply, err := c.client.Run(fullArgs...)

	if err != nil {
		return nil, fmt.Errorf("hotspot command failed: %w", err)
	}

	return reply, nil
}

// SetConfig sets client configuration
func (c *Client) SetConfig(config *Config) {
	c.config = config
}

// GetConfig returns client configuration
func (c *Client) GetConfig() *Config {
	return c.config
}
