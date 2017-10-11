package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/idmsubs/terraform-provider-scvmm/scvmm"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: scvmm.Provider})
}
