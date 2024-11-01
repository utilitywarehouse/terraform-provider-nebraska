package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccPackageResource_basic(t *testing.T) {
	dsn := "nebraska_package.test"

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourcePackage,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dsn, "id"),
					resource.TestCheckResourceAttrSet(dsn, "created_ts"),
					resource.TestCheckResourceAttrSet(dsn, "application_id"),
					resource.TestCheckResourceAttr(dsn, "version", "0.0.0"),
					resource.TestCheckResourceAttr(dsn, "arch", "amd64"),
					resource.TestCheckResourceAttr(dsn, "type", "flatcar"),
					resource.TestCheckResourceAttr(dsn, "url", "http://fake-address/"),
					resource.TestCheckResourceAttr(dsn, "filename", "test.tgz"),
					resource.TestCheckResourceAttr(dsn, "description", "Test package"),
					resource.TestCheckResourceAttr(dsn, "size", "465881871"),
					resource.TestCheckResourceAttr(dsn, "hash", "r3nufcxgMTZaxYEqL+x2zIoeClk="),
				),
			},
		},
	})
}

const testAccResourcePackage = `
provider "nebraska" {
}

resource "nebraska_package" "test_0" {
  version  = "0.1.0"
  arch     = "amd64"
  url      = "http://fake-address/"
}

resource "nebraska_channel" "test" {
  name       = "terraform-test"
  arch       = "amd64"
  package_id = nebraska_package.test_0.id
}

resource "nebraska_package" "test" {
  version  = "0.0.0"
  arch     = "amd64"
  url      = "http://fake-address/"
  type     = "flatcar"
  filename = "test.tgz"
  description = "Test package"
  size        = "465881871"
  hash        = "r3nufcxgMTZaxYEqL+x2zIoeClk="
  channels_blacklist = [
    nebraska_channel.test.id,
  ]
}
`
