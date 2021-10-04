package mssql

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"net/url"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password": {
				Type:     schema.TypeString,
				Required: true,
			},
			"host": {
				Type:     schema.TypeString,
				Required: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Default:  1433,
				Optional: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"mssql_database": resourceDatabase(),
			"mssql_user":     resourceUser(),
		},
		DataSourcesMap: map[string]*schema.Resource{},
		ConfigureContextFunc:  providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	u := &url.URL{
		Scheme: "sqlserver",
		User:   url.UserPassword(username, password),
		Host:   fmt.Sprintf("%s:%d", d.Get("host"), d.Get("port")),
	}

	db, err := sql.Open("sqlserver", u.String())
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return db, nil

}

