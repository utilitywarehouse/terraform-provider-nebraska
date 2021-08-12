package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccChannelResource_basic(t *testing.T) {
	dsn := "nebraska_channel.test"

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceChannel,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dsn, "id"),
					resource.TestCheckResourceAttr(dsn, "name", "test-terraform"),
					resource.TestCheckResourceAttr(dsn, "arch", "amd64"),
					resource.TestCheckResourceAttrSet(dsn, "package_id"),
					resource.TestCheckResourceAttr(dsn, "color", "#1fbb86"),
					resource.TestCheckResourceAttrSet(dsn, "created_ts"),
					resource.TestCheckResourceAttrSet(dsn, "application_id"),
				),
			},
		},
	})
}

const testAccResourceChannel = `
resource "nebraska_package" "test" {
  version = "0.0.0"
  arch    = "amd64"
  url     = "http://fake-address/"
}

resource "nebraska_channel" "test" {
  name       = "test-terraform"
  arch       = "amd64"
  package_id = nebraska_package.test.id
  color      = "#1fbb86"
}
`
