//go:generate go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@latest
//go:generate tfplugindocs generate --rendered-provider-name "SAP BTP"

package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"github.com/SAP/terraform-provider-btp/btp/provider"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	err := providerserver.Serve(context.Background(), provider.New, providerserver.ServeOpts{
		Address:         "registry.terraform.io/sap/btp",
		Debug:           debug,
		ProtocolVersion: 6,
	})

	if err != nil {
		log.Fatal(err)
	}
}
