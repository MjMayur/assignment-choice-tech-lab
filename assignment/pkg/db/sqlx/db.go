package sqlx

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const (
	MYSQL     = "mysql"
	MYSQLPORT = "3306"
)

type DbConfig struct {
	Driver, Host, User, Pass, Name string
}

type Database struct {
	ReadDB *sqlx.DB
}

func NewDBConn(dbConf *DbConfig) (*sqlx.DB, error) {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbConf.User, dbConf.Pass, dbConf.Host, MYSQLPORT, dbConf.Name)
	mysqlDbInstance, err := sqlx.Connect(dbConf.Driver, connectionString)
	if err != nil {
		return nil, err
	}
	return mysqlDbInstance, nil
}
