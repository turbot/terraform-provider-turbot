
resource "turbot_policy_setting" "test_policy" {
  resource = "arn:aws::eu-west-2:650022101893"
  policy_type = "tmod:@turbot/aws#/policy/types/accountStack"
  value = "Skip"
  precedence = "must"
}