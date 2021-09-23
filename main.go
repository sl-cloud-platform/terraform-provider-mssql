package main

import (
	"github.com/sl-cloud-platform/terraform-provider-mssql/mssql"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: mssql.Provider,
	})
}
