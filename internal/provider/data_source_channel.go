package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceChannel() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceChannelRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Description:  "Name of the channel.",
			},
			"arch": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"all", "amd64", "aarch64", "x86"}, false),
				Description:  "Arch.",
			},
			"application_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "ID of the application this channel belongs to.",
			},
			"color": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Hex color code of the channel on the UI.",
			},
			"created_ts": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation timestamp.",
			},
			"package_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of this channel's package.",
			},
		},
	}
}

func dataSourceChannelRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*apiClient)

	appID, err := getApplicationID(d, c)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("application_id", appID)
	channels, err := c.ListChannels(appID)
	if err != nil {
		return diag.FromErr(err)
	}
	name := d.Get("name").(string)
	arch := d.Get("arch").(string)

	for _, c := range channels {
		if c.Name == name && c.Arch.String() == arch {
			d.SetId(c.ID)
			d.Set("color", c.Color)
			d.Set("created_ts", c.CreatedTs.String())
			d.Set("package_id", c.PackageID)

			return nil
		}
	}

	return diag.Errorf("couldn't find channel %s (%s)", name, arch)
}
