package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func main() {
	db, err := sql.Open("mysql", "root@/plg")

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	rows, _ := db.Query("select * from components")

	for rows.Next() {
		var r1 int64
		var r2 string
		var r3 int
		var r4 string
		var r5 string
		if err := rows.Scan(&r1, &r2, &r3, &r4, &r5); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("plg: %d %s %d %s %s\n", r1, r2, r3, r4, r5)
	}
}
