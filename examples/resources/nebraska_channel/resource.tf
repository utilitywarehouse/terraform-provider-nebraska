data "nebraska_package" "package" {
  version = "2942.1.0"
  arch    = "amd64"
}

resource "nebraska_channel" "channel" {
  name       = "custom-channel"
  arch       = "amd64"
  package_id = data.nebraska_package.package.id
}
