
resource "turbot_folder" "parent" {
  parent = "tmod:@turbot/turbot#/"
  title = "provider_acceptance_tests"
  description = "Acceptance testing folder"
}
resource "turbot_policy_setting" "test_policy" {
  resource = turbot_folder.parent.id
  policy_type = "tmod:@turbot/provider-test#/policy/types/stringArrayPolicy"
  value = <<EOF
- a
- b
EOF
  precedence = "must"
}