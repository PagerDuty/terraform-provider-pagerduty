package main

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5/tf5server"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"

	"github.com/PagerDuty/terraform-provider-pagerduty/pagerduty"
	pagerdutyplugin "github.com/PagerDuty/terraform-provider-pagerduty/pagerdutyplugin"
)

func main() {
	Serve()
}

func Serve() {
	ctx := context.Background()

	muxServer, err := tf5muxserver.NewMuxServer(
		ctx,
		// terraform-plugin-framework
		providerserver.NewProtocol5(pagerdutyplugin.New()),
		// terraform-plugin-sdk
		pagerduty.Provider(pagerduty.IsMuxed).GRPCProvider,
	)
	if err != nil {
		log.Fatal(err)
	}

	var serveOpts []tf5server.ServeOpt

	address := "registry.terraform.io/pagerduty/pagerduty"
	err = tf5server.Serve(address, muxServer.ProviderServer, serveOpts...)
	if err != nil {
		log.Fatal(err)
	}
}
