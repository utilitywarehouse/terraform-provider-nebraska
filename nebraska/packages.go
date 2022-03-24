package nebraska

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/kinvolk/nebraska/backend/pkg/codegen"
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

	return ValidPackageTypes[i-1]
}

// PackageTypeFromString parses the string into a PackageType
func PackageTypeFromString(s string) (PackageType, error) {
	for i, sd := range ValidPackageTypes {
		if s == sd {
			return PackageType(i + 1), nil
		}

	}

	return PackageTypeOther, ErrInvalidPackageType
}

// GetPackage retrieves a package by its id
func (c *Client) GetPackage(appID, id string) (*codegen.Package, error) {
	req, err := c.newRequest(http.MethodGet, fmt.Sprintf("/api/apps/%s/packages/%s", appID, id), nil)
	if err != nil {
		return nil, err
	}

	data := &codegen.Package{}
	if err := c.do(req, data); err != nil {
		return nil, err
	}

	return data, nil
}

// SearchPackages lists the packages for a particular application and version
func (c *Client) SearchPackages(appID, id string) (*codegen.PackagePage, error) {
	req, err := c.newRequest(http.MethodGet, fmt.Sprintf("/api/apps/%s/packages?page=1&perpage=100000&searchVersion=%s", appID, id), nil)
	if err != nil {
		return nil, err
	}

	data := &codegen.PackagePage{}
	if err := c.do(req, data); err != nil {
		return nil, err
	}

	return data, nil
}

// AddPackageInput are the supported arguments when adding a package
type AddPackageInput struct {
	ApplicationID     string              `json:"application_id"`
	Arch              codegen.Arch        `json:"arch"`
	ChannelsBlacklist []string            `json:"channels_blacklist"`
	Description       string              `json:"description"`
	Filename          string              `json:"filename"`
	FlatcarAction     *FlatcarActionInput `json:"flatcar_action"`
	Hash              string              `json:"hash"`
	Size              string              `json:"size"`
	Type              PackageType         `json:"type"`
	URL               string              `json:"url"`
	Version           string              `json:"version"`
}

// FlatcarActionInput are the supported arguments when assigning a flatcar
// action to a package
type FlatcarActionInput struct {
	Sha256 string `json:"sha256"`
}

// AddPackage adds a new package
func (c *Client) AddPackage(appID string, input *AddPackageInput) (*codegen.Package, error) {
	req, err := c.newRequest(http.MethodPost, fmt.Sprintf("/api/apps/%s/packages", appID), input)
	if err != nil {
		return nil, err
	}

	data := &codegen.Package{}
	if err := c.do(req, data); err != nil {
		return nil, err
	}

	return data, nil
}

// UpdatePackageInput are the supported arguments when updating a package
type UpdatePackageInput struct {
	ApplicationID     string              `json:"application_id"`
	Arch              codegen.Arch        `json:"arch"`
	ChannelsBlacklist []string            `json:"channels_blacklist"`
	Description       string              `json:"description"`
	Filename          string              `json:"filename"`
	FlatcarAction     *FlatcarActionInput `json:"flatcar_action"`
	Hash              string              `json:"hash"`
	Size              string              `json:"size"`
	Type              PackageType         `json:"type"`
	URL               string              `json:"url"`
	Version           string              `json:"version"`
}

// UpdatePackage updates an existing package
func (c *Client) UpdatePackage(appID, id string, input *UpdatePackageInput) (*codegen.Package, error) {
	req, err := c.newRequest(http.MethodPut, fmt.Sprintf("/api/apps/%s/packages/%s", appID, id), input)
	if err != nil {
		return nil, err
	}

	data := &codegen.Package{}
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
