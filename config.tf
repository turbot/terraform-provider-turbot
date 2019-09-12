data "turbot_policy_value" "test_policy" {
  resource    = "tmod:@turbot/turbot#/"
  policy_type = "tmod:@turbot/turbot#/policy/types/domainName"
}
