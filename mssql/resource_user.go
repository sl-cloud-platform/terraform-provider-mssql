package mssql

import (
	"database/sql"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"errors"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserCreate,
		Read:   resourceUserRead,
		Update: resourceUserUpdate,
		Delete: resourceUserDelete,

		Schema: map[string]*schema.Schema{
			"database": {
				Type:     schema.TypeString,
				Required: true,
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceUserCreate(d *schema.ResourceData, m interface{}) error {
	db := m.(*sql.DB)
	database := d.Get("database").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	//ddladmin := d.Get("ddladmin").(bool)
	//datawriter := d.Get("datawriter").(bool)
	//datareader := d.Get("datareader").(bool)

	row, err := checkUser(db, username)
	// only try to create database if it not exists
	if err == sql.ErrNoRows {

		_, err := db.Query(fmt.Sprintf("CREATE LOGIN \"%s\" WITH PASSWORD = '%s', CHECK_POLICY = OFF, CHECK_EXPIRATION = OFF", username, password))
		if err != nil {
			return err
		}
		//create user ' + @user + ' for login ' + @login + ' with default_schema = dbo'
		//
		// TODO Schema?
		_, err = db.Query(fmt.Sprintf("exec('use %s; CREATE USER \"%s\" FOR LOGIN \"%s\" with default_schema = dbo')", database, username, username))
		//_, err = db.Query(fmt.Sprintf(  "CREATE USER \"%s\" FOR LOGIN \"%s\" with default_schema = dbo", username, username))
		if err != nil {
			return err
		}
	}

	row, err = checkUser(db, username)
	if err != nil {
		return err
	}
	d.SetId(fmt.Sprint(row.principal_id))

	return err

	//if err = &row.principal_id; err != nil {
	//	return err
	//}
	//
	//d.SetId(fmt.Sprint(id))
	//return err

	//row, err = checkTable(db, name)
	//if err != nil {
	//	return err
	//}
	//d.SetId(row.name)
	//
	//return err

	//// add permissions
	//if ddladmin {
	//	//_, err = db.Query(fmt.Sprintf(  "CREATE USER '%s' FOR LOGIN '%s' with default_schema = dbo", username, username))
	//	//exec('use ' + @db + '; alter role db_ddladmin add member ' + @user)
	//	_, err = db.Query(fmt.Sprintf(  "exec('use '%s'; alter role db_ddladmin add member '%s'", db, username))
	//	if err != nil {
	//		return err
	//	}
	//}
	//
	//if datawriter {
	//	//_, err = db.Query(fmt.Sprintf(  "CREATE USER '%s' FOR LOGIN '%s' with default_schema = dbo", username, username))
	//	//exec('use ' + @db + '; alter role db_ddladmin add member ' + @user)
	//	_, err = db.Query(fmt.Sprintf(  "exec('use '%s'; alter role db_datawriter add member '%s'", db, username))
	//	if err != nil {
	//		return err
	//	}
	//}
	//
	//if datareader {
	//	//_, err = db.Query(fmt.Sprintf(  "CREATE USER '%s' FOR LOGIN '%s' with default_schema = dbo", username, username))
	//	//exec('use ' + @db + '; alter role db_ddladmin add member ' + @user)
	//	_, err = db.Query(fmt.Sprintf(  "exec('use '%s'; alter role db_datareader add member '%s'", db, username))
	//	if err != nil {
	//		return err
	//	}
	//}


	//row := db.QueryRow(fmt.Sprintf("SELECT principal_id FROM master.sys.server_principals WHERE username = '%s'", username))
	//var id int
	//if err = row.Scan(&id); err != nil {
	//	return err
	//}
	//
	//d.SetId(fmt.Sprint(id))
	//return err
}

type PrinicipalsRow struct {
	principal_id int
}

func checkUser(db *sql.DB, username string) (*PrinicipalsRow, error) {

	var row PrinicipalsRow
	err := db.QueryRow(fmt.Sprintf("SELECT principal_id FROM master.sys.server_principals where name = '%s'", username)).Scan(&row.principal_id)
	if err != nil {
		return nil, errors.New(fmt.Sprint("check user", err))
	}
	return &row, nil
}

func resourceUserRead(d *schema.ResourceData, m interface{}) error {

	db := m.(*sql.DB)
	row := db.QueryRow(fmt.Sprintf("SELECT name FROM master.sys.server_principals WHERE principal_id = %s", d.Id()))
	var name string
	err := row.Scan(&name)
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		return err
	}
	if err := d.Set("username", name); err != nil {
		return err
	}
	return nil
}

func resourceUserUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceUserDelete(d *schema.ResourceData, m interface{}) error {
	db := m.(*sql.DB)
	row := db.QueryRow(fmt.Sprintf("SELECT name FROM master.sys.server_principals WHERE principal_id = %s", d.Id()))
	var name string
	err := row.Scan(&name)

	if err != sql.ErrNoRows {
		dropUserQuery := fmt.Sprintf("DROP USER %s", name)
		_, err = db.Query(dropUserQuery)
		if err != nil {
			return errors.New(fmt.Sprint("Failed to drop user", err))
		}
		_, err = db.Query(fmt.Sprintf("DROP LOGIN %s", name))
		if err != nil {
			return errors.New(fmt.Sprint("Failed to drop login", err))
		}
	}

	return nil
}
