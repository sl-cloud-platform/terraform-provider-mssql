# Terraform Microsoft SQL Server Provider

## Usage
```hcl
provider "mssql" {
  host = "localhost"
  username = "sa"
  password = "password"
}

resource "mssql_database" "db" {
  name = "MyDatabase"
}

resource "mssql_user" "user" {
  database = mssql_database.db.name 
  name = "MyUser"
  password = "MyPassword"
}
```