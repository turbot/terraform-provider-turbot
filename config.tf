
resource "turbot_policy_setting" "test_policy" {
  resource = "tmod:@turbot/turbot#/"
  policy_type = "tmod:@turbot/aws#/policy/types/accountStack"
  value = "Skip"
}
