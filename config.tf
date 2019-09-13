resource "turbot_mod" "test" {
  parent  = "tmod:@turbot/turbot#/"
  org     = "turbot"
  mod     = "provider-test"
  version = "*"
}
