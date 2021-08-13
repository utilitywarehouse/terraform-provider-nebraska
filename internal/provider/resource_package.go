package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/utilitywarehouse/terraform-provider-nebraska/nebraska"
)

func resourcePackage() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePackageCreate,
		ReadContext:   resourcePackageRead,
		UpdateContext: resourcePackageUpdate,
		DeleteContext: resourcePackageDelete,

		Schema: map[string]*schema.Schema{
			"version": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Description:  "Package version.",
			},
			"url": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsURLWithHTTPorHTTPS,
				Description:  "URL where the package is available.",
			},
			"arch": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(nebraska.ValidArchs, false),
				Default:      nebraska.ArchAll.String(),
				Description:  "Package arch.",
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(nebraska.ValidPackageTypes, false),
				Default:      nebraska.PackageTypeFlatcar.String(),
				Description:  "Type of package.",
			},
			"filename": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The filename of the package.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A description of the package.",
			},
			"size": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The size, in bytes.",
			},
			"hash": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A base64 encoded sha1 hash of the package digest. Tip: `cat update.gz | openssl dgst -sha1 -binary | base64`.",
			},
			"channels_blacklist": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "A list of channels (by id) that cannot point to this package.",
			},
			"flatcar_action": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
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
							Type:        schema.TypeString,
							Required:    true,
							Description: "A base64 encoded sha256 hash of the action. Tip: `cat update.gz | openssl dgst -sha256 -binary | base64`.",
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
			"application_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "ID of the application this package belongs to.",
			},
			"created_ts": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation timestamp.",
			},
		},
	}
}

func resourcePackageCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*apiClient)

	appID, err := getApplicationID(d, c)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("application_id", appID)

	arch, err := nebraska.ArchFromString(d.Get("arch").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	pkgType, err := nebraska.PackageTypeFromString(d.Get("type").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	input := &nebraska.AddPackageInput{
		Type:              pkgType,
		Version:           d.Get("version").(string),
		URL:               d.Get("url").(string),
		Filename:          d.Get("filename").(string),
		Description:       d.Get("description").(string),
		Size:              d.Get("size").(string),
		Hash:              d.Get("hash").(string),
		ChannelsBlacklist: expandChannelBlacklist(d.Get("channels_blacklist").([]interface{})),
		Arch:              arch,
		FlatcarAction: &nebraska.FlatcarActionInput{
			Sha256: expandFlatcarActionSha256(d.Get("flatcar_action").([]interface{})),
		},
	}

	pkg, err := c.AddPackage(appID, input)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(pkg.ID)

	return resourcePackageRead(ctx, d, meta)
}

func resourcePackageRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*apiClient)

	appID, err := getApplicationID(d, c)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("application_id", appID)

	pkg, err := c.GetPackage(appID, d.Id())
	if err != nil {
		if err == nebraska.ErrNotFound {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	if pkg == nil {
		d.SetId("")
		return nil
	}

	if err := d.Set("type", pkg.Type.String()); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("version", pkg.Version); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("url", pkg.URL); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("filename", pkg.Filename); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("description", pkg.Description); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("size", pkg.Size); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("created_ts", pkg.CreatedTs.String()); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("channels_blacklist", pkg.ChannelsBlacklist); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("flatcar_action", flattenFlatcarAction(pkg.FlatcarAction)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("arch", pkg.Arch.String()); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourcePackageUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*apiClient)

	appID, err := getApplicationID(d, c)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("application_id", appID)

	arch, err := nebraska.ArchFromString(d.Get("arch").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	pkgType, err := nebraska.PackageTypeFromString(d.Get("type").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	input := &nebraska.UpdatePackageInput{
		Type:              pkgType,
		Version:           d.Get("version").(string),
		URL:               d.Get("url").(string),
		Filename:          d.Get("filename").(string),
		Description:       d.Get("description").(string),
		Size:              d.Get("size").(string),
		Hash:              d.Get("hash").(string),
		ChannelsBlacklist: expandChannelBlacklist(d.Get("channels_blacklist").([]interface{})),
		Arch:              arch,
		FlatcarAction: &nebraska.FlatcarActionInput{
			Sha256: expandFlatcarActionSha256(d.Get("flatcar_action").([]interface{})),
		},
	}

	if _, err := c.UpdatePackage(appID, d.Id(), input); err != nil {
		return diag.FromErr(err)
	}

	return resourcePackageRead(ctx, d, meta)
}

func resourcePackageDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*apiClient)

	appID, err := getApplicationID(d, c)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := c.DeletePackage(appID, d.Id()); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func expandChannelBlacklist(l []interface{}) []string {
	var blacklist []string
	for _, i := range l {
		if v, ok := i.(string); ok {
			blacklist = append(blacklist, v)
		}
	}

	return blacklist
}

func expandFlatcarActionSha256(l []interface{}) string {
	if len(l) == 0 || l[0] == nil {
		return ""
	}

	m := l[0].(map[string]interface{})

	if v, ok := m["sha256"].(string); ok {
		return v
	}

	return ""
}

func flattenFlatcarAction(action *nebraska.FlatcarAction) []map[string]interface{} {
	if action == nil {
		return []map[string]interface{}{}
	}

	return []map[string]interface{}{
		{
			"id":                      action.ID,
			"event":                   action.Event,
			"chromeos_version":        action.ChromeOSVersion,
			"sha256":                  action.Sha256,
			"needs_admin":             action.NeedsAdmin,
			"is_delta":                action.IsDelta,
			"disable_payload_backoff": action.DisablePayloadBackoff,
			"metadata_signature_rsa":  action.MetadataSignatureRsa,
			"metadata_size":           action.MetadataSize,
			"deadline":                action.Deadline,
			"created_ts":              action.CreatedTs.String(),
		},
	}
}
