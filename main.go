package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6/tf6server"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"

	"github.com/mackerelio-labs/terraform-provider-mackerel/mackerel"
	"github.com/mackerelio-labs/terraform-provider-mackerel/mackerel-v1/provider"
)

const (
	providerAddr = "registry.terraform.io/mackerelio-labs/mackerel"
)

func main() {
	var debug bool
	flag.BoolVar(&debug, "debug", false, "run as debug-mode")

	flag.Parse()

	ctx := context.Background()

	upgradedSDKProvider, err := tf5to6server.UpgradeServer(
		ctx,
		mackerel.Provider().GRPCProvider,
	)
	if err != nil {
		log.Fatal(err)
	}

	providers := []func() tfprotov6.ProviderServer{
		func() tfprotov6.ProviderServer {
			return upgradedSDKProvider
		},
		providerserver.NewProtocol6(provider.New()),
	}

	muxServer, err := tf6muxserver.NewMuxServer(ctx, providers...)
	if err != nil {
		log.Fatal(err)
	}

	var serveOpts []tf6server.ServeOpt
	if debug {
		serveOpts = append(serveOpts, tf6server.WithManagedDebug())
	}

	err = tf6server.Serve(
		providerAddr,
		muxServer.ProviderServer,
		serveOpts...,
	)
	if err != nil {
		log.Fatal(err)
	}
}
