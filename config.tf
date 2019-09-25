resource "turbot_google_directory" "test" {
	title = "google_directory_test_provider2"
	profile_id_template = "profileemail"
	status = "New"
	directory_type = "google"
	client_id = "GoogleDirTest4"
	client_secret = "fb-tbevaACsBKQHthzba-PH9"
	parent = "tmod:@turbot/turbot#/"
	description = "test Directory"
}
}