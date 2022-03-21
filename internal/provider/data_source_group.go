package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceGroup() *schema.Resource {
	return &schema.Resource{
		Description: "A group provides a particular release channel to machines and controls various options that manage the update procedure.",
		ReadContext: dataSourceGroupRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Description:  "Name of the group.",
			},
			"application_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "ID of the application this group belongs to.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A description of the group.",
			},
			"created_ts": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation timestamp.",
			},
			"rollout_in_progress": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether a rollout is currently in progress for this group.",
			},
			"channel_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The channel this group provides.",
			},
			"policy_updates_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Are updates enabled?",
			},
			"policy_safe_mode": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Safe mode will only update 1 instance at a time, and stop if an update fails.",
			},
			"policy_office_hours": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Only update between 9am and 5pm.",
			},
			"policy_timezone": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Timezone used to inform `policy_office_hours`.",
			},
			"policy_period_interval": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Period used in combination with `policy_max_updates_per_period`.",
			},
			"policy_max_updates_per_period": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The maximum number of updates that can be performed within the `policy_period_interval`.",
			},
			"policy_update_timeout": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Timeout for updates.",
			},
			"track": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Identifier for clients.",
			},
		},
	}
}

func dataSourceGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*apiClient)

	appID, err := getApplicationID(d, c)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("application_id", appID)
	groupPage, err := c.ListGroups(appID)
	if err != nil {
		return diag.FromErr(err)
	}
	if groupPage.Count != groupPage.TotalCount {
		return diag.FromErr(fmt.Errorf("GET groups returned %d/%d groups. We don't paginate.", groupPage.Count, groupPage.TotalCount))
	}
	name := d.Get("name").(string)
	for _, g := range groupPage.Groups {
		if g.Name == name {
			d.SetId(g.Id)
			d.Set("description", g.Description)
			d.Set("created_ts", g.CreatedTs.String())
			d.Set("rollout_in_progress", g.RolloutInProgress)
			d.Set("channel_id", g.ChannelID)
			d.Set("policy_updates_enabled", g.PolicyUpdatesEnabled)
			d.Set("policy_safe_mode", g.PolicySafeMode)
			d.Set("policy_office_hours", g.PolicyOfficeHours)
			d.Set("policy_timezone", g.PolicyTimezone)
			d.Set("policy_period_interval", g.PolicyPeriodInterval)
			d.Set("policy_max_updates_per_period", g.PolicyMaxUpdatesPerPeriod)
			d.Set("policy_update_timeout", g.PolicyUpdateTimeout)
			d.Set("track", g.Track)

			return nil
		}
	}

	return diag.Errorf("couldn't find group")
}
