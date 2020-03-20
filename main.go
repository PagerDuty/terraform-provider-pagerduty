package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/mLupine/terraform-provider-pagerduty/pagerduty"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: pagerduty.Provider})
}
