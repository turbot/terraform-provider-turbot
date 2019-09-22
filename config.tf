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
	parent = "tmod:@turbot/turbot#/"
	title = "${lookup(var.user_details, "${element(keys(var.user_details), count.index)}")}"
	display_name = "${lookup(var.user_details, "${element(keys(var.user_details), count.index)}")}"
	email = "${element(keys(var.user_details), count.index)}"
	given_name = "${element(split(" ", lookup(var.user_details, "${element(keys(var.user_details), count.index)}")), 0)}"
	family_name = "${element(split(" ", lookup(var.user_details, "${element(keys(var.user_details), count.index)}")), 1)}"
	directory_pool_id = "pool1"
	status = "Active"
	profile_id = "${element(keys(var.user_details), count.index)}"

}