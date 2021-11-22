#!/bin/bash

set -x

provider_version=0.2.3
provider_path=registry.terraform.io/sl-cloud-platform/mssql/"$provider_version"/darwin_amd64/

go build -o terraform-provider-mssql_"$provider_version"

#mkdir -p ~/Library/Application\ Support/io.terraform/plugins/"$provider_path"
#cp terraform-provider-mssql_"$provider_version"  ~/Library/Application\ Support/io.terraform/plugins/"$provider_path"