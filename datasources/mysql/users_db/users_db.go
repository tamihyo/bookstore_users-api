package users_db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
 
	_ "github.com/go-sql-driver/mysql"
)

// const (
// 	mysql_users_username = "mysql_users_username"
// 	mysql_users_password = "mysql_users_password"
// 	mysql_users_host     = "mysql_users_host"
// 	mysql_users_port     = "mysql_users_port"
// 	mysql_users_schema   = "mysql_users_schema"
// )

var (
	Client *sql.DB

	// username = os.Getenv("MYSQL_USERNAME")
	// password = os.Getenv("mysql_users_password")
	// host     = os.Getenv("mysql_users_host")
	// port     = os.Getenv("mysql_users_port")
	// schema   = os.Getenv("mysql_users_schema")
)

func InitDB() {
	username := os.Getenv("mysql_users_username")
	password := os.Getenv("mysql_users_password")
	host := os.Getenv("mysql_users_host")
	port := os.Getenv("mysql_users_port")
	schema := os.Getenv("mysql_users_schema")

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8",
		username,
		password,
		host,
		port,
		schema,
	)
	/*
		That's exactly what we should pay attention to.
		if we're working on a method or
		function the variables declared with := are scoped to
		that particular function. If we're trying to create the
		global Client variable then we
		should avoid the := operator since we're
		creating a local variable in that case.
	*/
	/*declare variable error to make global Client variable*/

	var err error

	Client, err = sql.Open("mysql", dataSourceName)

	if err != nil {
		panic(err)
	}

	if err = Client.Ping(); err != nil {
		panic(err)
	}

	log.Println("database success configured", Client)

}
