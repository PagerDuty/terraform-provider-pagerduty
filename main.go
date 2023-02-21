package main

import (
	"github.com/arqiva-tb/terraform-provider-pagerduty/pagerduty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: pagerduty.Provider})
}
