package model

import (
	"database/sql"
)

type AssignmentUser struct {
	ID          int            `db:"id"`
	FirstName   string         `db:"first_name"`
	LastName    string         `db:"last_name"`
	CompanyName sql.NullString `db:"company_name"`
	Address     sql.NullString `db:"address"`
	City        sql.NullString `db:"city"`
	County      sql.NullString `db:"county"`
	Postal      sql.NullString `db:"postal"`
	Phone       sql.NullString `db:"phone"`
	Email       sql.NullString `db:"email"`
	Web         sql.NullString `db:"web"`
	CreatedAt   sql.NullString `db:"created_at"`
	UpdatedAt   sql.NullString `db:"updated_at"`
	DeletedAt   sql.NullString `db:"deleted_at"`
}

var AssignmentUserModelMap = map[string]FieldStruct{
	"ID":          {MySQLDatatype: "int", FieldName: "id"},
	"FirstName":   {MySQLDatatype: "varchar", FieldName: "first_name"},
	"LastName":    {MySQLDatatype: "varchar", FieldName: "last_name"},
	"CompanyName": {MySQLDatatype: "varchar", FieldName: "company_name"},
	"Address":     {MySQLDatatype: "varchar", FieldName: "address"},
	"City":        {MySQLDatatype: "varchar", FieldName: "city"},
	"County":      {MySQLDatatype: "varchar", FieldName: "county"},
	"Postal":      {MySQLDatatype: "varchar", FieldName: "postal"},
	"Phone":       {MySQLDatatype: "varchar", FieldName: "phone"},
	"Email":       {MySQLDatatype: "varchar", FieldName: "email"},
	"Web":         {MySQLDatatype: "varchar", FieldName: "web"},
	"CreatedAt":   {MySQLDatatype: "timestamp", FieldName: "created_at"},
	"UpdatedAt":   {MySQLDatatype: "timestamp", FieldName: "updated_at"},
	"DeletedAt":   {MySQLDatatype: "timestamp", FieldName: "deleted_at"},
}
