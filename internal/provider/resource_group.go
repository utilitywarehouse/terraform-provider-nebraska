package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/utilitywarehouse/terraform-provider-nebraska/nebraska"
)

func resourceGroup() *schema.Resource {
	return &schema.Resource{
		Description: "A group provides a particular release channel to machines and controls various options that manage the update procedure.",

		CreateContext: resourceGroupCreate,
		ReadContext:   resourceGroupRead,
		UpdateContext: resourceGroupUpdate,
		DeleteContext: resourceGroupDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Description:  "Name of the group.",
			},
			"track": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Identifier for clients, filled with the group ID if omitted.",
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
				Optional:    true,
				Description: "A description of the group.",
			},
			"created_ts": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation timestamp",
			},
			"rollout_in_progress": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether a rollout is currently in progress for this group.",
			},
			"channel_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The channel this group provides.",
			},
			"policy_updates_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Enable updates.",
			},
			"policy_safe_mode": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Safe mode will only update 1 instance at a time, and stop if an update fails.",
			},
			"policy_office_hours": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Only update between 9am and 5pm.",
			},
			"policy_timezone": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Timezone used to inform `policy_office_hours`.",
			},
			"policy_period_interval": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "1 minutes",
				Description: "Period used in combination with `policy_max_updates_per_period`.",
			},
			"policy_max_updates_per_period": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     9999999,
				Description: "The maximum number of updates that can be performed within the `policy_period_interval`.",
			},
			"policy_update_timeout": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "60 minutes",
				Description: "Timeout for updates",
			},
		},
	}
}

func resourceGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*apiClient)

	appID, err := getApplicationID(d, c)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("application_id", appID)

	input := &nebraska.AddGroupInput{
		Name:                      d.Get("name").(string),
		Description:               d.Get("description").(string),
		ChannelID:                 d.Get("channel_id").(string),
		PolicyUpdatesEnabled:      d.Get("policy_updates_enabled").(bool),
		PolicySafeMode:            d.Get("policy_safe_mode").(bool),
		PolicyOfficeHours:         d.Get("policy_office_hours").(bool),
		PolicyTimezone:            d.Get("policy_timezone").(string),
		PolicyPeriodInterval:      d.Get("policy_period_interval").(string),
		PolicyMaxUpdatesPerPeriod: d.Get("policy_max_updates_per_period").(int),
		PolicyUpdateTimeout:       d.Get("policy_update_timeout").(string),
		Track:                     d.Get("track").(string),
	}

	group, err := c.AddGroup(appID, input)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(group.Id)

	return resourceGroupRead(ctx, d, meta)
}

func resourceGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*apiClient)

	appID, err := getApplicationID(d, c)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("application_id", appID)
	group, err := c.GetGroup(appID, d.Id())
	if err != nil {
		if err == nebraska.ErrNotFound {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	if group == nil {
		d.SetId("")
		return nil
	}

	if err := d.Set("name", group.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("description", group.Description); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("created_ts", group.CreatedTs.String()); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("rollout_in_progress", group.RolloutInProgress); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("channel_id", group.ChannelID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("policy_updates_enabled", group.PolicyUpdatesEnabled); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("policy_safe_mode", group.PolicySafeMode); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("policy_office_hours", group.PolicyOfficeHours); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("policy_timezone", group.PolicyTimezone); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("policy_period_interval", group.PolicyPeriodInterval); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("policy_max_updates_per_period", group.PolicyMaxUpdatesPerPeriod); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("policy_update_timeout", group.PolicyUpdateTimeout); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("track", group.Track); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*apiClient)

	appID, err := getApplicationID(d, c)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("application_id", appID)

	input := &nebraska.UpdateGroupInput{
		Name:                      d.Get("name").(string),
		Description:               d.Get("description").(string),
		ChannelID:                 d.Get("channel_id").(string),
		PolicyUpdatesEnabled:      d.Get("policy_updates_enabled").(bool),
		PolicySafeMode:            d.Get("policy_safe_mode").(bool),
		PolicyOfficeHours:         d.Get("policy_office_hours").(bool),
		PolicyTimezone:            d.Get("policy_timezone").(string),
		PolicyPeriodInterval:      d.Get("policy_period_interval").(string),
		PolicyMaxUpdatesPerPeriod: d.Get("policy_max_updates_per_period").(int),
		PolicyUpdateTimeout:       d.Get("policy_update_timeout").(string),
		Track:                     d.Get("track").(string),
	}

	if _, err := c.UpdateGroup(appID, d.Id(), input); err != nil {
		return diag.FromErr(err)
	}

	return resourceGroupRead(ctx, d, meta)
}

func resourceGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*apiClient)

	appID, err := getApplicationID(d, c)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := c.DeleteGroup(appID, d.Id()); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
