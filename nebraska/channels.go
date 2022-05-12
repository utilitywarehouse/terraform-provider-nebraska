package nebraska

import (
	"fmt"
	"net/http"

	"github.com/kinvolk/nebraska/backend/pkg/codegen"
)

// GetChannel retrieves a channel by its id
func (c *Client) GetChannel(appID, channelID string) (*codegen.Channel, error) {
	req, err := c.newRequest(http.MethodGet, fmt.Sprintf("/api/apps/%s/channels/%s", appID, channelID), nil)
	if err != nil {
		return nil, err
	}

	data := &codegen.Channel{}
	if err := c.do(req, data); err != nil {
		return nil, err
	}

	return data, nil
}

// ListChannels lists the channels for a particular application
func (c *Client) ListChannels(appID string) (*codegen.ChannelPage, error) {
	req, err := c.newRequest(http.MethodGet, fmt.Sprintf("/api/apps/%s/channels?page=1&perpage=10000", appID), nil)
	if err != nil {
		return nil, err
	}

	data := &codegen.ChannelPage{}
	if err := c.do(req, data); err != nil {
		return nil, err
	}

	return data, nil
}

// AddChannelInput are the supported arguments when adding a channel
type AddChannelInput struct {
	Name          string       `json:"name"`
	Color         string       `json:"color"`
	PackageID     string       `json:"package_id"`
	ApplicationID string       `json:"application_id"`
	Arch          codegen.Arch `json:"arch"`
}

// AddChannel adds a new channel
func (c *Client) AddChannel(appID string, input *AddChannelInput) (*codegen.Channel, error) {
	req, err := c.newRequest(http.MethodPost, fmt.Sprintf("/api/apps/%s/channels", appID), input)
	if err != nil {
		return nil, err
	}

	data := &codegen.Channel{}
	if err := c.do(req, data); err != nil {
		return nil, err
	}

	return data, nil
}

// UpdateChannelInput are the supported arguments when updating a channel
type UpdateChannelInput struct {
	Name          string       `json:"name"`
	Color         string       `json:"color"`
	PackageID     string       `json:"package_id"`
	ApplicationID string       `json:"application_id"`
	Arch          codegen.Arch `json:"arch"`
}

// UpdateChannel updates an existing channel
func (c *Client) UpdateChannel(appID, id string, input *UpdateChannelInput) (*codegen.Channel, error) {
	req, err := c.newRequest(http.MethodPut, fmt.Sprintf("/api/apps/%s/channels/%s", appID, id), input)
	if err != nil {
		return nil, err
	}

	data := &codegen.Channel{}
	if err := c.do(req, data); err != nil {
		return nil, err
	}

	return data, nil
}

// DeleteChannel removes a channel
func (c *Client) DeleteChannel(appID, id string) error {
	req, err := c.newRequest(http.MethodDelete, fmt.Sprintf("/api/apps/%s/channels/%s", appID, id), nil)
	if err != nil {
		return err
	}

	return c.do(req, nil)
}
