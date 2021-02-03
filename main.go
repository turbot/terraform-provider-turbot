package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/turbot/terraform-provider-turbot/turbot"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: turbot.Provider})
}
