package mssql

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			"roles": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
		},
	}
}

func resourceUserCreate(d *schema.ResourceData, m interface{}) error {
	db := m.(*sql.DB)
	database := d.Get("database").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	roles := d.Get("roles").(*schema.Set).List()

	row, err := checkUser(db, username)
	if err == sql.ErrNoRows {

		_, err := db.Query(fmt.Sprintf("CREATE LOGIN \"%s\" WITH PASSWORD = '%s', CHECK_POLICY = OFF, CHECK_EXPIRATION = OFF", username, password))
		if err != nil {
			return errors.New(fmt.Sprint("Failed to create login", err))
		}

		// TODO Schema?
		_, err = db.Query(fmt.Sprintf("exec('use %s; CREATE USER \"%s\" FOR LOGIN \"%s\" with default_schema = dbo')", database, username, username))
		//_, err = db.Query(fmt.Sprintf(  "CREATE USER \"%s\" FOR LOGIN \"%s\" with default_schema = dbo", username, username))
		if err != nil {
			return errors.New(fmt.Sprint("Failed to create user", err))
		}

	}

	row, err = checkUser(db, username)
	if err != nil {
		return errors.New(fmt.Sprint("Unknow error occured", err))
	}

	for _, role := range roles {
		_, err = db.Exec(fmt.Sprintf("exec('use %s; exec(''sp_addrolemember %s, %s '') ')", database, role, username))
		if err != nil {
			return errors.New(fmt.Sprint("Failed to add member to role ", err))
		}
	}

	d.SetId(fmt.Sprint(row.principal_id))

	return err
}

type PrinicipalsRow struct {
	principal_id int
}

func checkUser(db *sql.DB, username string) (*PrinicipalsRow, error) {

	var row PrinicipalsRow
	err := db.QueryRow(fmt.Sprintf("SELECT principal_id FROM master.sys.server_principals where name = '%s'", username)).Scan(&row.principal_id)
	if err != nil {
		return nil, err
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
	database := d.Get("database").(string)
	row := db.QueryRow(fmt.Sprintf("SELECT name FROM master.sys.server_principals WHERE principal_id = %s", d.Id()))
	var name string
	err := row.Scan(&name)

	if err != sql.ErrNoRows {
		_, err = db.Query(fmt.Sprintf("DROP LOGIN %s", name))
		if err != nil {
			return errors.New(fmt.Sprint("Failed to drop login", err))
		}
		_, err = db.Query(fmt.Sprintf("exec('use %s; drop user %s');", database, name))
		if err != nil {
			return errors.New(fmt.Sprint("Failed to drop user", err))
		}
	}

	return nil
}
