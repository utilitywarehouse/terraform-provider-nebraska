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

