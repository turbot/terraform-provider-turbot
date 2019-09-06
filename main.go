package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/terraform-providers/terraform-provider-turbot/turbot"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: turbot.Provider})
}
