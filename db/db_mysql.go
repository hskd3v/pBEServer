package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/harriklein/pBE/pBEServer/log"

	// force call
	_ "github.com/go-sql-driver/mysql"
)

const (
	envDBpbeMysqlUsername = "DB_PBE_MYSQL_USERNAME"
	envDBpbeMysqlPassword = "DB_PBE_MYSQL_PASSWORD"
	envDBpbeMysqlHost     = "DB_PBE_MYSQL_HOST"
	envDBpbeMysqlSchema   = "DB_PBE_MYSQL_SCHEMA"
)

var (
	// ConnPBE is a global
	ConnPBE *sql.DB

	username = os.Getenv(envDBpbeMysqlUsername)
	password = os.Getenv(envDBpbeMysqlPassword)
	host     = os.Getenv(envDBpbeMysqlHost)
	schema   = os.Getenv(envDBpbeMysqlSchema)
)

// Init initializes the database
var mysqlInit = func() {

	log.Log.Println("Connecting into DB (MySQL)...")

	_dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",
		username, password, host, schema,
	)

	var _error error
	ConnPBE, _error = sql.Open("mysql", _dataSourceName)
	if _error != nil {
		panic(_error)
	}

	if _error = ConnPBE.Ping(); _error != nil {
		panic(_error)
	}

	log.Log.Debugln("Database sucessfully configured")
}
