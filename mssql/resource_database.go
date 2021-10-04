package mssql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDatabase() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDatabaseCreate,
		ReadContext:   resourceDatabaseRead,
		UpdateContext: resourceDatabaseUpdate,
		DeleteContext: resourceDatabaseDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"drop_on_destroy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceDatabaseCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	db := m.(*sql.DB)
	name := d.Get("name").(string)

	row, err := checkDatabase(db, name)
	// only try to create database if it not exists
	if err == sql.ErrNoRows {
		_, err := db.Query(fmt.Sprintf("CREATE DATABASE %s", name))
		if err != nil {
			return diag.FromErr(err)
		}
	}
	row, err = checkDatabase(db, name)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(row.name)

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceDatabaseRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	db := m.(*sql.DB)
	name := d.Id()
	row, err := checkDatabase(db, name)
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("name", row.name); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceDatabaseUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func resourceDatabaseDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	dropOnDestroy := d.Get("drop_on_destroy").(bool)
	name := d.Id()

	//return errors.New(name)

	// TODO fix drop database, raises this error sometimes: Error: mssql: Warning: Fatal error 615 occurred at Sep 23 2021 12:18PM. Note the error and time, and contact your system administrator.
	//USE [master]
	//GO
	//ALTER DATABASE [Beratungsmappe] SET SINGLE_USER WITH ROLLBACK IMMEDIATE
	//GO
	//USE [master]
	//GO
	//DROP DATABASE [Beratungsmappe]
	//GO

	if dropOnDestroy {
		db := m.(*sql.DB)
		row, err := checkDatabase(db, name)
		if err != nil && err != sql.ErrNoRows {
			return diag.FromErr(errors.New(fmt.Sprint("Failed to check if database exists", err)))
		}

		if row != nil {

			_, err := db.Exec(fmt.Sprintf("DROP DATABASE %s", name))
			if err != nil {
				return diag.FromErr(errors.New(fmt.Sprint("Failed to drop database", err)))
			}
		}
	}

	return nil
}
