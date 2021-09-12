package todb

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

func ConnectDB() *sql.DB {
	cfg := mysql.Config{
		User:                 os.Getenv("DBUSER"), //NOTE: PLEASE SET THIS ENV IN YOUR MACHINE
		Passwd:               os.Getenv("DBPASS"), //NOTE: PLEASE SET THIS ENV IN YOUR MACHINE
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "btvn_b11",
		AllowNativePasswords: true,
	}
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	if pingErr := db.Ping(); pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println(db, *db)
	fmt.Println("Connected")
	return db
}
