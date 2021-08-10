package nebraska

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

// PackageType is the type of package
type PackageType int

const (
	// PackageTypeFlatcar is a Flatcar update package
	PackageTypeFlatcar PackageType = 1 + iota
	// PackageTypeDocker is a docker container
	PackageTypeDocker
	// PackageTypeRocket is a rkt container
	PackageTypeRocket
	// PackageTypeOther is a generic package type
	PackageTypeOther
)

var (
	// ErrInvalidPackageType is a custom error returned when an unsupported arch is
	// requested
	ErrInvalidPackageType = errors.New("nebraska: invalid/unsupported package type")

	// ValidPackageTypes are the package types that Nebraska supports
	ValidPackageTypes = []string{
		"flatcar",
		"docker",
		"rkt",
		"other",
	}
)

// String returns the string representation of the package type
func (pt PackageType) String() string {
	i := int(pt)

	return ValidPackageTypes[i]
}

// PackageTypeFromString parses the string into a PackageType
func PackageTypeFromString(s string) (PackageType, error) {
	for i, sd := range ValidPackageTypes {
		if s == sd {
			return PackageType(i), nil
		}

	}

	return PackageTypeOther, ErrInvalidPackageType
}

// Package is a package in Nebraska
type Package struct {
	ID                string         `json:"id"`
	Type              PackageType    `json:"type"`
	Version           string         `json:"version"`
	URL               string         `json:"url"`
	Filename          string         `json:"filename"`
	Description       string         `json:"description"`
	Size              string         `json:"size"`
	Hash              string         `json:"hash"`
	CreatedTs         time.Time      `json:"created_ts"`
	ChannelsBlacklist []string       `json:"channels_blacklist"`
	ApplicationID     string         `json:"application_id"`
	FlatcarAction     *FlatcarAction `json:"flatcar_action"`
	Arch              Arch           `json:"arch"`
}

// FlatcarAction is a flatcar action
type FlatcarAction struct {
	ID                    string    `json:"id"`
	Event                 string    `json:"event"`
	ChromeOSVersion       string    `json:"chromeos_version"`
	Sha256                string    `json:"sha256"`
	NeedsAdmin            bool      `json:"needs_admin"`
	IsDelta               bool      `json:"is_delta"`
	DisablePayloadBackoff bool      `json:"disable_payload_backoff"`
	MetadataSignatureRsa  string    `json:"metadata_signature_rsa"`
	MetadataSize          string    `json:"metadata_size"`
	Deadline              string    `json:"deadline"`
	CreatedTs             time.Time `json:"created_ts"`
}

// GetPackage retrieves a package by its id
func (c *Client) GetPackage(appID, id string) (*Package, error) {
	req, err := c.newRequest(http.MethodGet, fmt.Sprintf("/api/apps/%s/packages/%s", appID, id), nil)
	if err != nil {
		return nil, err
	}

	data := &Package{}
	if err := c.do(req, data); err != nil {
		return nil, err
	}

	return data, nil
}

// ListPackages lists the packages for a particular application
func (c *Client) ListPackages(appID string) ([]*Package, error) {
	req, err := c.newRequest(http.MethodGet, fmt.Sprintf("/api/apps/%s/packages", appID), nil)
	if err != nil {
		return nil, err
	}

	data := []*Package{}
	if err := c.do(req, &data); err != nil {
		return nil, err
	}

	return data, nil
}

// AddPackageInput are the supported arguments when adding a package
type AddPackageInput struct {
	Type              PackageType         `json:"type"`
	Version           string              `json:"version"`
	URL               string              `json:"url"`
	Filename          string              `json:"filename"`
	Description       string              `json:"description"`
	Size              string              `json:"size"`
	Hash              string              `json:"hash"`
	ChannelsBlacklist []string            `json:"channels_blacklist"`
	Arch              Arch                `json:"arch"`
	FlatcarAction     *FlatcarActionInput `json:"flatcar_action"`
}

// FlatcarActionInput are the supported arguments when assigning a flatcar
// action to a package
type FlatcarActionInput struct {
	Sha256 string `json:"sha256"`
}

// AddPackage adds a new package
func (c *Client) AddPackage(appID string, input *AddPackageInput) (*Package, error) {
	req, err := c.newRequest(http.MethodPost, fmt.Sprintf("/api/apps/%s/packages", appID), input)
	if err != nil {
		return nil, err
	}

	data := &Package{}
	if err := c.do(req, data); err != nil {
		return nil, err
	}

	return data, nil
}

// UpdatePackageInput are the supported arguments when updating a package
type UpdatePackageInput struct {
	Type              PackageType         `json:"type"`
	Version           string              `json:"version"`
	URL               string              `json:"url"`
	Filename          string              `json:"filename"`
	Description       string              `json:"description"`
	Size              string              `json:"size"`
	Hash              string              `json:"hash"`
	ChannelsBlacklist []string            `json:"channels_blacklist"`
	Arch              Arch                `json:"arch"`
	FlatcarAction     *FlatcarActionInput `json:"flatcar_action"`
}

// UpdatePackage updates an existing package
func (c *Client) UpdatePackage(appID, id string, input *UpdatePackageInput) (*Package, error) {
	req, err := c.newRequest(http.MethodPut, fmt.Sprintf("/api/apps/%s/packages/%s", appID, id), input)
	if err != nil {
		return nil, err
	}

	data := &Package{}
	if err := c.do(req, data); err != nil {
		return nil, err
	}

	return data, nil
}

// DeletePackage deletes a package
func (c *Client) DeletePackage(appID, id string) error {
	req, err := c.newRequest(http.MethodDelete, fmt.Sprintf("/api/apps/%s/packages/%s", appID, id), nil)
	if err != nil {
		return err
	}

	return c.do(req, nil)
}
