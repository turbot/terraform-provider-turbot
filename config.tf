################################        Variables        ###################################

variable "local_directory_name" {
	type = "string"
	default = "test_directory"
}

variable "user_details" {
	type = "map"
	default = {
		# email = "Display Name"
		"test_user1@turbot.com" = "Test User1"
		"test_user2@turbot.com" = "Test User2"
	}
}

################################        Resources        ###################################

resource "turbot_local_directory" "test_dir" {
	parent = "tmod:@turbot/turbot#/"
	title = "${var.local_directory_name}"
	description = "test Directory"
	profile_id_template = "{{profile.email}}"
}

resource "turbot_local_directory_user" "test_user" {
	count = "${length(var.user_details)}"
	title = "${lookup(var.user_details, "${element(keys(var.user_details), count.index)}")}"
	email = "${element(keys(var.user_details), count.index)}"
	status = "Active"
	display_name = "${lookup(var.user_details, "${element(keys(var.user_details), count.index)}")}"
	parent = "${turbot_local_directory.test_dir.id}"
}

resource "turbot_profile" "test_user" {
	count = "${length(var.user_details)}"
	title = "${lookup(var.user_details, "${element(keys(var.user_details), count.index)}")}"
	email = "${element(keys(var.user_details), count.index)}"
	status = "Active"
	directory_pool_id = "dpi"
	given_name = "${element(split(" ", lookup(var.user_details, "${element(keys(var.user_details), count.index)}")), 0)}"
	family_name = "${element(split(" ", lookup(var.user_details, "${element(keys(var.user_details), count.index)}")), 1)}"
	display_name = "${lookup(var.user_details, "${element(keys(var.user_details), count.index)}")}"
	parent = "${turbot_local_directory.test_dir.id}"
	status = "Active"
	profile_id = "${element(keys(var.user_details), count.index)}"
}

resource "turbot_grant" "test" {
	count = "${length(var.user_details)}"
	resource = "tmod:@turbot/turbot#/"
	permission_type = "tmod:@turbot/aws#/permission/types/aws"
	permission_level = "tmod:@turbot/turbot-iam#/permission/levels/superuser"
	profile_id = "${turbot_profile.test_user.*.id[count.index]}"
}