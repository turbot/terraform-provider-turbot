data "turbot_policy_value" "test_policy" {
  resource    = "tmod:@turbot/turbot#/"
  policy_type = "tmod:@turbot/turbot#/policy/types/domainName"
}


output "op1"{
  value = "${data.turbot_policy_value.test_policy.value}"
}