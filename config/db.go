package config

import (
	"database/sql"
	"fmt"
	"log"
)

var Database *sql.DB

func ConnectDB() {
	var err error
	Database, err = sql.Open("mysql", "root:Learn@321@tcp(localhost)/learnathon")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("DB Connected")
	// defer Database.Close()
}
