package main

import (
	"context"
	"log"

	pd "github.com/PagerDuty/terraform-provider-pagerduty/pagerduty"
	pdp "github.com/PagerDuty/terraform-provider-pagerduty/pagerdutyplugin"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6/tf6server"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
)

func main() {
	Serve()
}

func Serve() {
	ctx := context.Background()

	upgradedSdkServer, err := tf5to6server.UpgradeServer(ctx, pd.Provider(pd.IsMuxed).GRPCProvider)
	if err != nil {
		log.Fatal(err)
	}

	muxServer, err := tf6muxserver.NewMuxServer(
		ctx,
		providerserver.NewProtocol6(pdp.New()),
		func() tfprotov6.ProviderServer { return upgradedSdkServer },
	)
	if err != nil {
		log.Fatal(err)
	}

	address := "registry.terraform.io/pagerduty/pagerduty"
	if err != nil {
		log.Fatal(err)
	}

	err = tf6server.Serve(address, muxServer.ProviderServer)
	if err != nil {
		log.Fatal(err)
	}
}
