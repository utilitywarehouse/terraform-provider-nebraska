package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccGroupDataSource_basic(t *testing.T) {
	dsn := "data.nebraska_group.test"

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGroup,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dsn, "id"),
					resource.TestCheckResourceAttrSet(dsn, "rollout_in_progress"),
					resource.TestCheckResourceAttrSet(dsn, "channel_id"),
					resource.TestCheckResourceAttrSet(dsn, "created_ts"),
					resource.TestCheckResourceAttrSet(dsn, "application_id"),
					resource.TestCheckResourceAttr(dsn, "name", "terraform-test"),
					resource.TestCheckResourceAttr(dsn, "description", "Test description"),
					resource.TestCheckResourceAttr(dsn, "policy_updates_enabled", "false"),
					resource.TestCheckResourceAttr(dsn, "policy_safe_mode", "true"),
					resource.TestCheckResourceAttr(dsn, "policy_office_hours", "true"),
					resource.TestCheckResourceAttr(dsn, "policy_timezone", "Europe/Berlin"),
					resource.TestCheckResourceAttr(dsn, "policy_period_interval", "10 minutes"),
					resource.TestCheckResourceAttr(dsn, "policy_max_updates_per_period", "10"),
					resource.TestCheckResourceAttr(dsn, "policy_update_timeout", "35 minutes"),
				),
			},
		},
	})
}

const testAccDataSourceGroup = `
resource "nebraska_package" "test" {
  version = "0.0.0"
  arch    = "amd64"
  url     = "http://fake-address/"
}

resource "nebraska_channel" "test" {
  name       = "terraform-test"
  arch       = "amd64"
  package_id = nebraska_package.test.id
}

resource "nebraska_group" "test" {
  name                           = "terraform-test"
  track                          = "terraform-test"
  description                    = "Test description"
  channel_id                     = nebraska_channel.test.id
  policy_updates_enabled         = "false"
  policy_safe_mode               = "true"
  policy_office_hours            = "true"
  policy_timezone                = "Europe/Berlin"
  policy_period_interval         = "10 minutes"
  policy_max_updates_per_period  = 10
  policy_update_timeout          = "35 minutes"
}

data "nebraska_group" "test" {
 name = nebraska_group.test.name
}
`
