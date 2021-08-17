package nebraska

import (
	"fmt"
	"net/http"
	"time"
)

// Group is a group in Nebraska
type Group struct {
	ID                        string    `json:"id"`
	Name                      string    `json:"name"`
	Description               string    `json:"description"`
	CreatedTs                 time.Time `json:"created_ts"`
	RolloutInProgress         bool      `json:"rollout_in_progress"`
	ApplicationID             string    `json:"application_id"`
	ChannelID                 string    `json:"channel_id"`
	PolicyUpdatesEnabled      bool      `json:"policy_updates_enabled"`
	PolicySafeMode            bool      `json:"policy_safe_mode"`
	PolicyOfficeHours         bool      `json:"policy_office_hours"`
	PolicyTimezone            string    `json:"policy_timezone"`
	PolicyPeriodInterval      string    `json:"policy_period_interval"`
	PolicyMaxUpdatesPerPeriod int       `json:"policy_max_updates_per_period"`
	PolicyUpdateTimeout       string    `json:"policy_update_timeout"`
	Track                     string    `json:"track"`
}

// GetGroup retrieves a group by its id
func (c *Client) GetGroup(appID, id string) (*Group, error) {
	req, err := c.newRequest(http.MethodGet, fmt.Sprintf("/api/apps/%s/groups/%s", appID, id), nil)
	if err != nil {
		return nil, err
	}

	data := &Group{}
	if err := c.do(req, data); err != nil {
		return nil, err
	}

	return data, nil
}

// ListGroups lists the groups for a particular application
func (c *Client) ListGroups(appID string) ([]*Group, error) {
	req, err := c.newRequest(http.MethodGet, fmt.Sprintf("/api/apps/%s/groups", appID), nil)
	if err != nil {
		return nil, err
	}

	data := []*Group{}
	if err := c.do(req, &data); err != nil {
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
func (c *Client) AddGroup(appID string, input *AddGroupInput) (*Group, error) {
	req, err := c.newRequest(http.MethodPost, fmt.Sprintf("/api/apps/%s/groups", appID), input)
	if err != nil {
		return nil, err
	}

	data := &Group{}
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
func (c *Client) UpdateGroup(appID, id string, input *UpdateGroupInput) (*Group, error) {
	req, err := c.newRequest(http.MethodPut, fmt.Sprintf("/api/apps/%s/groups/%s", appID, id), input)
	if err != nil {
		return nil, err
	}

	data := &Group{}
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
