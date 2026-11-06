package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/turbot/terraform-provider-turbot/turbot"
)

// version is set by the release build (-X main.version).
var version = "dev"

func main() {
	turbot.ProviderVersion = version
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: turbot.Provider})
}
