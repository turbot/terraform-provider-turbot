resource "turbot_mod" "test" {
  parent = "tmod:@turbot/turbot#/"
  org = "turbot"
  mod = "structure-test"
  version = "5.0.0-beta.59"
}

resource "turbot_folder" "parent" {
  parent = "tmod:@turbot/turbot#/"
  title = "provider_test_parent"
  description = "PARENT"
  //  This is a workaround for odd mod install issue
  // without it mod install will fail with `Runnable already has a process running` error for modInstall control
  depends_on = [
    "turbot_mod.test"]

}

resource "turbot_folder" "child" {
  parent = turbot_folder.parent.id
  title = "provider_test_child"
  description = "CHILD"
}

resource "turbot_policy_setting" "test_policy" {
  resource = turbot_folder.child.id
  policy_type = "${turbot_mod.test.uri}#/policy/types/testPolicy"
  value = "TEST"
  precedence = "should"
}

