resource "turbot_policy_setting" "test_policy" {
  resource       = "arn:aws::eu-west-2:650022101893"
  policy_type    = "tmod:@turbot/aws#/policy/types/accountStack"
  template_input = "{ account{ Id } }"
  template       = "{% if $.account.Id == '650022101893' %}Skip{% else %}'Check: Configured'{% endif %}"
  precedence     = "should"
}
