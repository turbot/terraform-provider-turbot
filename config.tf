resource "turbot_resource" "test" {
  parent = "tmod:@turbot/turbot#/"
  type = "tmod:@turbot/turbot#/resource/types/folder"
  payload =
  <<





  EOF
  {


    "title": "provider_test2",


    "description": "test resource2"


  }
  EOF
}