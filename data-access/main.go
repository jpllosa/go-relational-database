package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
	// Capture connection properties.
	// c:\> set DBUSER=golang cmd promt
	// $ export DBUSER=golang linux/mac
	// $env:DBUSER = "golang"

	cfg := mysql.Config{
		// 	User:   os.Getenv("DBUSER"),
		// 	Passwd: os.Getenv("DBPASS"),
		User:   "golang",
		Passwd: "golang",
		Net:    "tcp",
		// Addr:                 "127.0.0.1:3306",
		Addr:                 "localhost:3306",
		DBName:               "recordings",
		AllowNativePasswords: true, // if not included, mysql native password authentication error is generated
	}
	// cfg := mysql.NewConfig()
	// cfg.User = os.Getenv("DBUSER")
	// cfg.Passwd = os.Getenv("DBPASS")
	// cfg.User = "golang"
	// cfg.Passwd = "golang"
	// cfg.Net = "tcp"
	// cfg.Addr = "127.0.0.1:3306" // works on mysql-5.7.33, workbench from host localhost, https://dba.stackexchange.com/questions/38803/mysql-error-1045-28000-access-denied-for-user
	// cfg.Addr = "localhost:3306" // works on mysql-5.6.43, workbench from host localhost
	// cfg.DBName = "recordings"
	// Get a database handle.
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")
}
