package nebraska

import (
	"fmt"
	"net/http"

	"github.com/kinvolk/nebraska/backend/pkg/codegen"
)

// GetGroup retrieves a group by its id
func (c *Client) GetGroup(appID, id string) (*codegen.Group, error) {
	req, err := c.newRequest(http.MethodGet, fmt.Sprintf("/api/apps/%s/groups/%s", appID, id), nil)
	if err != nil {
		return nil, err
	}

	data := &codegen.Group{}
	if err := c.do(req, data); err != nil {
		return nil, err
	}

	return data, nil
}

// ListGroups lists the groups for a particular application
func (c *Client) ListGroups(appID string) (*codegen.GroupPage, error) {
	req, err := c.newRequest(http.MethodGet, fmt.Sprintf("/api/apps/%s/groups?page=1&perpage=10000", appID), nil)
	if err != nil {
		return nil, err
	}

	data := &codegen.GroupPage{}
	if err := c.do(req, data); err != nil {
		return nil, err
	}

	return data, nil
}

// AddGroupInput are the supported arguments when adding a channel
type AddGroupInput struct {
	Name                      string `json:"name"`
	Description               string `json:"description"`
	ChannelID                 string `json:"channel_id"`
	PolicyUpdatesEnabled      bool   `json:"policy_updates_enabled"`
	PolicySafeMode            bool   `json:"policy_safe_mode"`
	PolicyOfficeHours         bool   `json:"policy_office_hours"`
	PolicyTimezone            string `json:"policy_timezone"`
	PolicyPeriodInterval      string `json:"policy_period_interval"`
	PolicyMaxUpdatesPerPeriod int    `json:"policy_max_updates_per_period"`
	PolicyUpdateTimeout       string `json:"policy_update_timeout"`
	Track                     string `json:"track"`
}

// AddGroup adds a new group
func (c *Client) AddGroup(appID string, input *AddGroupInput) (*codegen.Group, error) {
	req, err := c.newRequest(http.MethodPost, fmt.Sprintf("/api/apps/%s/groups", appID), input)
	if err != nil {
		return nil, err
	}

	data := &codegen.Group{}
	if err := c.do(req, data); err != nil {
		return nil, err
	}

	return data, nil
}

// UpdateGroupInput are the supported arguments when updating a channel
type UpdateGroupInput struct {
	Name                      string `json:"name"`
	Description               string `json:"description"`
	ChannelID                 string `json:"channel_id"`
	PolicyUpdatesEnabled      bool   `json:"policy_updates_enabled"`
	PolicySafeMode            bool   `json:"policy_safe_mode"`
	PolicyOfficeHours         bool   `json:"policy_office_hours"`
	PolicyTimezone            string `json:"policy_timezone"`
	PolicyPeriodInterval      string `json:"policy_period_interval"`
	PolicyMaxUpdatesPerPeriod int    `json:"policy_max_updates_per_period"`
	PolicyUpdateTimeout       string `json:"policy_update_timeout"`
	Track                     string `json:"track"`
}

// UpdateGroup updates an existing group
func (c *Client) UpdateGroup(appID, id string, input *UpdateGroupInput) (*codegen.Group, error) {
	req, err := c.newRequest(http.MethodPut, fmt.Sprintf("/api/apps/%s/groups/%s", appID, id), input)
	if err != nil {
		return nil, err
	}

	data := &codegen.Group{}
	if err := c.do(req, data); err != nil {
		return nil, err
	}

	return data, nil
}

// DeleteGroup deletes a group
func (c *Client) DeleteGroup(appID, id string) error {
	req, err := c.newRequest(http.MethodDelete, fmt.Sprintf("/api/apps/%s/groups/%s", appID, id), nil)
	if err != nil {
		return err
	}

	return c.do(req, nil)
}
