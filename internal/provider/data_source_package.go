package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourcePackage() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePackageRead,
		Schema: map[string]*schema.Schema{
			"version": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"arch": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"all", "amd64", "aarch64", "x86"}, false),
			},
			"application_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"filename": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"hash": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_ts": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"flatcar_action": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Computed: true,
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
	packages, err := c.ListPackages(appID)
	if err != nil {
		return diag.FromErr(err)
	}
	version := d.Get("version").(string)
	arch := d.Get("arch").(string)

	for _, p := range packages {
		if p.Version == version && p.Arch.String() == arch {
			d.SetId(p.ID)
			d.Set("type", p.Type.String())
			d.Set("url", p.URL)
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
