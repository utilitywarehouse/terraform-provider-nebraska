package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/kinvolk/nebraska/backend/pkg/api"
	"github.com/utilitywarehouse/terraform-provider-nebraska/nebraska"
)

func dataSourcePackage() *schema.Resource {
	return &schema.Resource{
		Description: "A versioned package of the application.",
		ReadContext: dataSourcePackageRead,
		Schema: map[string]*schema.Schema{
			"version": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Description:  "Package version.",
			},
			"arch": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"all", "amd64", "aarch64", "x86"}, false),
				Description:  "Package arch.",
			},
			"application_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of package.",
			},
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL where the package is available.",
			},
			"filename": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The filename of the package.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A description of the package.",
			},
			"size": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The size, in bytes.",
			},
			"hash": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A base64 encoded sha1 hash of the package digest.",
			},
			"created_ts": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation timestamp.",
			},
			"flatcar_action": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A Flatcar specific Omaha action.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"event": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"chromeos_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sha256": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"needs_admin": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_delta": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"disable_payload_backoff": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"metadata_signature_rsa": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"metadata_size": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"deadline": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_ts": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"channels_blacklist": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "A list of channels (by id) that cannot point to this package.",
			},
		},
	}
}

func dataSourcePackageRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*apiClient)

	appID, err := getApplicationID(d, c)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("application_id", appID)

	version := d.Get("version").(string)
	arch := d.Get("arch").(string)

	packagePage, err := c.SearchPackages(appID, version)
	if err != nil {
		return diag.FromErr(err)
	}
	if packagePage.Count != packagePage.TotalCount {
		return diag.FromErr(fmt.Errorf("GET packages returned %d/%d packages. We don't paginate.", packagePage.Count, packagePage.TotalCount))
	}

	for _, p := range packagePage.Packages {
		if p.Version == version && api.Arch(p.Arch).String() == arch {
			d.SetId(p.Id)
			d.Set("type", nebraska.PackageType(p.Type).String())
			d.Set("url", p.Url)
			d.Set("filename", p.Filename)
			d.Set("description", p.Description)
			d.Set("size", p.Size)
			d.Set("hash", p.Hash)
			d.Set("created_ts", p.CreatedTs.String())
			d.Set("channels_blacklist", p.ChannelsBlacklist)
			d.Set("flatcar_action", flattenFlatcarAction(p.FlatcarAction))

			return nil
		}
	}

	return diag.Errorf("couldn't find package")
}
