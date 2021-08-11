package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGroupRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"application_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_ts": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rollout_in_progress": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"channel_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"policy_updates_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"policy_safe_mode": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"policy_office_hours": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"policy_timezone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"policy_period_interval": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"policy_max_updates_per_period": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"policy_update_timeout": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"track": {
				Type:     schema.TypeString,
				Computed: true,
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
	groups, err := c.ListGroups(appID)
	if err != nil {
		return diag.FromErr(err)
	}
	name := d.Get("name").(string)
	for _, g := range groups {
		if g.Name == name {
			d.SetId(g.ID)
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
