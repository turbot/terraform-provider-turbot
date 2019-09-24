resource "turbot_smart_folder" "test" {
		parent  = 				"tmod:@turbot/turbot#/"
		filters = 			["arn:aws:iam::013122550996:user/pratik/accesskey/AKIAQGDRKHTKBON32K3J"],
		description =     "Smart Folder Testing"
		title = 					"smart_folder"
	}