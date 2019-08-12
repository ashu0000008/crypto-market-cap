package db

import (
	"database/sql"
	"github.com/ashu0000008/crypto-market-cap/db/config"
	_ "github.com/go-sql-driver/mysql"
)

var Database *sql.DB

// Setup the database connection
func Setup() (*sql.DB, error) {

	dbUser := config.GetString("database.user")
	dbPasswd := config.GetString("database.passwd")
	dbHost := config.GetString("database.host")
	dbName := config.GetString("database.name")
	dbConnection := config.GetString("database.connection")
	connectionString := dbUser + ":" + dbPasswd + "@" + dbConnection + "(" + dbHost + ")" + "/" + dbName + "?charset=utf8"

	db, err := sql.Open("mysql", connectionString)

	if err != nil {
		print(err)
		return nil, err
	}

	// Ping the database once since Open() doesn't open a connection
	err = db.Ping()
	if err != nil {
		print(err)
		return nil, err
	}

	Database = db
	return Database, nil
}
