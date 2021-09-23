package mssql

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDatabase() *schema.Resource {
	return &schema.Resource{
		Create: resourceDatabaseCreate,
		Read:   resourceDatabaseRead,
		Update: resourceDatabaseUpdate,
		Delete: resourceDatabaseDelete,

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

func resourceDatabaseCreate(d *schema.ResourceData, m interface{}) error {
	db := m.(*sql.DB)
	name := d.Get("name").(string)

	row, err := checkDatabase(db, name)
	// only try to create database if it not exists
	if err == sql.ErrNoRows {
		_, err := db.Query(fmt.Sprintf("CREATE DATABASE %s", name))
		if err != nil {
			return err
		}
	}
	row, err = checkDatabase(db, name)
	if err != nil {
		return err
	}
	d.SetId(row.name)

	return err
}

type DatabaseSchemaRow struct {
	name string
}

func resourceDatabaseRead(d *schema.ResourceData, m interface{}) error {
	db := m.(*sql.DB)
	name := d.Id()
	row, err := checkDatabase(db, name)
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		return err
	}

	if err = d.Set("name", row.name); err != nil {
		return err
	}

	return nil
}

func checkDatabase(db *sql.DB, name string) (*DatabaseSchemaRow, error) {
	var row DatabaseSchemaRow
	err := db.QueryRow(fmt.Sprintf("SELECT name FROM sys.databases where name = '%s'", name)).Scan(&row.name)
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func resourceDatabaseUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceDatabaseDelete(d *schema.ResourceData, m interface{}) error {
	dropOnDestroy := d.Get("drop_on_destroy").(bool)
	name := d.Id()

	//return errors.New(name)

	// TODO fix drop database, raises this error sometimes: Error: mssql: Warning: Fatal error 615 occurred at Sep 23 2021 12:18PM. Note the error and time, and contact your system administrator.

	if dropOnDestroy {
		db := m.(*sql.DB)
		row, err := checkDatabase(db, name)
		if err != nil && err != sql.ErrNoRows {
			return errors.New(fmt.Sprint("Failed to check if database exists", err))
		}

		if row != nil {
			_, err := db.Exec(fmt.Sprintf("DROP DATABASE %s", name))
			if err != nil {
				return errors.New(fmt.Sprint("Failed to drop database", err))
			}
		}
	}

	return nil
}
