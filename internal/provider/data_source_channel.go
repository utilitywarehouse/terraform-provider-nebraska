package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/kinvolk/nebraska/backend/pkg/api"
)

func dataSourceChannel() *schema.Resource {
	return &schema.Resource{
		Description: "A release channel that provides a particular package version.",
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
	channelPage, err := c.ListChannels(appID)
	if err != nil {
		return diag.FromErr(err)
	}
	if channelPage.Count != channelPage.TotalCount {
		return diag.FromErr(fmt.Errorf("GET channels returned %d/%d channels. We don't paginate.", channelPage.Count, channelPage.TotalCount))
	}
	name := d.Get("name").(string)
	arch := d.Get("arch").(string)

	for _, c := range channelPage.Channels {
		if c.Name == name && api.Arch(c.Arch).String() == arch {
			d.SetId(c.Id)
			d.Set("color", c.Color)
			d.Set("created_ts", c.CreatedTs.String())
			d.Set("package_id", c.PackageID)

			return nil
		}
	}

	return diag.Errorf("couldn't find channel %s (%s)", name, arch)
}
