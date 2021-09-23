# Terraform provider for Microsoft SQL Server at AWS

## Usage
```hcl
provider "mssql" {
  host = "localhost"
  username = "sa"
  password = "password"
}

resource "mssql_database" "db" {
  name = "MyDatabase"
  drop_on_destroy = true
}

resource "mssql_user" "user" {
  database = mssql_database.db.name 
  name = "MyUser"
  password = "MyPassword"
  roles = ["db_ddladmin", "db_datawriter", "db_datareader"]
}
```