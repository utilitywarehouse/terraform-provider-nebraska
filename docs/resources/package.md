---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "nebraska_package Resource - terraform-provider-nebraska"
subcategory: ""
description: |-
  A versioned package of the application.
---

# nebraska_package (Resource)

A versioned package of the application.

## Example Usage

```terraform
resource "nebraska_package" "test" {
  type        = "flatcar"
  version     = "2191.5.0"
  url         = "https://update.release.flatcar-linux.net/amd64-usr/2191.5.0/"
  filename    = "flatcar_production_update.gz"
  description = "Flatcar Container Linux 2191.5.0"
  size        = "465881871"
  hash        = "r3nufcxgMTZaxYEqL+x2zIoeClk="

  flatcar_action {
    sha256 = "LIkAKVZY2EJFiwTmltiJZLFLA5xT/FodbjVgqkyF/y8="
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `url` (String) URL where the package is available.
- `version` (String) Package version.

### Optional

- `application_id` (String) ID of the application this package belongs to.
- `arch` (String) Package arch. Defaults to `all`.
- `channels_blacklist` (List of String) A list of channels (by id) that cannot point to this package.
- `description` (String) A description of the package.
- `filename` (String) The filename of the package.
- `flatcar_action` (Block List, Max: 1) A Flatcar specific Omaha action. (see [below for nested schema](#nestedblock--flatcar_action))
- `hash` (String) A base64 encoded sha1 hash of the package digest. Tip: `cat update.gz | openssl dgst -sha1 -binary | base64`.
- `size` (String) The size, in bytes.
- `type` (String) Type of package. Defaults to `flatcar`.

### Read-Only

- `created_ts` (String) Creation timestamp.
- `id` (String) The ID of this resource.

<a id="nestedblock--flatcar_action"></a>
### Nested Schema for `flatcar_action`

Required:

- `sha256` (String) A base64 encoded sha256 hash of the action. Tip: `cat update.gz | openssl dgst -sha256 -binary | base64`.

Read-Only:

- `chromeos_version` (String)
- `created_ts` (String)
- `deadline` (String)
- `disable_payload_backoff` (Boolean)
- `event` (String)
- `id` (String)
- `is_delta` (Boolean)
- `metadata_signature_rsa` (String)
- `metadata_size` (String)
- `needs_admin` (Boolean)
