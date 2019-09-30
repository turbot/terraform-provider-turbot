package turbot

//
//import (
//	"fmt"
//	"github.com/hashicorp/terraform/helper/resource"
//	"github.com/hashicorp/terraform/terraform"
//	"github.com/terraform-providers/terraform-provider-turbot/apiclient"
//	"testing"
//)
//
//// test suites
//// TODO these fail currently - awaiting a mod update - turbot is currently a required property
//func TestAccResourceAwsAccount(t *testing.T) {
//	resource.Test(t, resource.TestCase{
//		PreCheck:     func() { testAccPreCheck(t) },
//		Providers:    testAccProviders,
//		CheckDestroy: testAccCheckResourceDestroy,
//		Steps: []resource.TestStep{
//			{
//				Config: testAccResourceConfig(awsAccountType, awsAccountPayload),
//				Check: resource.ComposeTestCheckFunc(
//					testAccCheckResourceExists("turbot_resource.test"),
//					resource.TestCheckResourceAttr(
//						"turbot_resource.test", "type", awsAccountType),
//					resource.TestCheckResourceAttr(
//						"turbot_resource.test", "payload", formatPayload(awsAccountPayload)),
//				),
//			},
//			{
//				Config: testAccResourceConfig(awsAccountType, awsAccountPayloadUpdateId),
//				Check: resource.ComposeTestCheckFunc(
//					testAccCheckResourceExists("turbot_resource.test"),
//					resource.TestCheckResourceAttr(
//						"turbot_resource.test", "type", awsAccountType),
//					resource.TestCheckResourceAttr(
//						"turbot_resource.test", "payload", formatPayload(awsAccountPayloadUpdateId)),
//				),
//			},
//			{
//				Config: testAccResourceConfig(awsAccountType, awsAccountPayloadUpdateEmail),
//				Check: resource.ComposeTestCheckFunc(
//					testAccCheckResourceExists("turbot_resource.test"),
//					resource.TestCheckResourceAttr(
//						"turbot_resource.test", "type", awsAccountType),
//					resource.TestCheckResourceAttr(
//						"turbot_resource.test", "payload", formatPayload(awsAccountPayloadUpdateEmail)),
//				),
//			},
//			{
//				Config: testAccResourceConfig(awsAccountType, awsAccountPayloadUpdateName),
//				Check: resource.ComposeTestCheckFunc(
//					testAccCheckResourceExists("turbot_resource.test"),
//					resource.TestCheckResourceAttr(
//						"turbot_resource.test", "type", awsAccountType),
//					resource.TestCheckResourceAttr(
//						"turbot_resource.test", "payload", formatPayload(awsAccountPayloadUpdateName)),
//				),
//			},
//		},
//	})
//}
//
//func TestAccResourceFolder(t *testing.T) {
//	resource.Test(t, resource.TestCase{
//		PreCheck:     func() { testAccPreCheck(t) },
//		Providers:    testAccProviders,
//		CheckDestroy: testAccCheckResourceDestroy,
//		Steps: []resource.TestStep{
//			{
//				Config: testAccResourceConfig(folderType, folderPayload),
//				Check: resource.ComposeTestCheckFunc(
//					testAccCheckResourceExists("turbot_resource.test"),
//					resource.TestCheckResourceAttr(
//						"turbot_resource.test", "type", folderType),
//					resource.TestCheckResourceAttr(
//						"turbot_resource.test", "payload", formatPayload(folderPayload)),
//				),
//			},
//
//			{
//				Config: testAccResourceConfig(folderType, folderPayloadUpdatedDescription),
//				Check: resource.ComposeTestCheckFunc(
//					testAccCheckResourceExists("turbot_resource.test"),
//					resource.TestCheckResourceAttr(
//						"turbot_resource.test", "type", folderType),
//					//
//					resource.TestCheckResourceAttr(
//						"turbot_resource.test", "payload", formatPayload(folderPayloadUpdatedDescription)),
//				),
//			},
//			{
//				Config: testAccResourceConfig(folderType, folderPayloadUpdatedTitle),
//				Check: resource.ComposeTestCheckFunc(
//					testAccCheckResourceExists("turbot_resource.test"),
//					resource.TestCheckResourceAttr(
//						"turbot_resource.test", "type", folderType),
//					resource.TestCheckResourceAttr(
//						"turbot_resource.test", "payload", formatPayload(folderPayloadUpdatedTitle)),
//				),
//			},
//		},
//	})
//}
//
//// configs
//var folderType = `tmod:@turbot/turbot#/resource/types/folder`
//var folderPayload = `{
//  "title": "provider_test",
//  "description": "test resource"
//}
//`
//var folderPayloadUpdatedTitle = `{
//  "title": "provider_test_",
//  "description": "test resource"
//}
//`
//var folderPayloadUpdatedDescription = `{
//  "title": "provider_test_",
//  "description": "test resource"
//}
//`
//
//var awsAccountType = `tmod:@turbot/aws#/resource/types/account`
//var awsAccountPayload = `{
// "Id": "123456789999"
//}
//`
//var awsAccountPayloadUpdateId = `{
// "Id": "123456781111"
//}
//`
//var awsAccountPayloadUpdateEmail = `{
// "Id": "123456789999",
// "Email": "kai@turbot.com"
//}
//`
//var awsAccountPayloadUpdateName = `{
// "Id": "123456789999",
// "Name": "kai"
//}
//`
//
//func testAccResourceConfig(resourceType, payload string) string {
//	return fmt.Sprintf(`
//resource "turbot_resource" "test" {
//  parent = "tmod:@turbot/turbot#/"
//  type = "%s"
//  payload =  <<EOF
//%sEOF
//}
//`, resourceType, payload)
//}
//
//// helper functions
//func testAccCheckResourceExists(resource string) resource.TestCheckFunc {
//	return func(state *terraform.State) error {
//		rs, ok := state.RootModule().Resources[resource]
//		if !ok {
//			return fmt.Errorf("Not found: %s", resource)
//		}
//		if rs.Primary.ID == "" {
//			return fmt.Errorf("No Record ID is set")
//		}
//		client := testAccProvider.Meta().(*apiclient.Client)
//		_, err := client.ReadResource(rs.Primary.ID, nil)
//		if err != nil {
//			return fmt.Errorf("error fetching item with resource %s. %s", resource, err)
//		}
//		return nil
//	}
//}
//
//func testAccCheckResourceDestroy(s *terraform.State) error {
//	client := testAccProvider.Meta().(*apiclient.Client)
//	for _, rs := range s.RootModule().Resources {
//		if rs.Type != "resource" {
//			continue
//		}
//		_, err := client.ReadResource(rs.Primary.ID, nil)
//		if err == nil {
//			return fmt.Errorf("Alert still exists")
//		}
//		if !apiclient.NotFoundError(err) {
//			return fmt.Errorf("expected 'not found' error, got %s", err)
//		}
//	}
//
//	return nil
//}
