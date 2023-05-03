package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var Db *sql.DB

func InitialiseDb() {

	fmt.Println("Initialising Db!!")
	Db, err := sql.Open("sqlite3", "./users.db")
	if err != nil {
		log.Fatal(err)
	}
	defer Db.Close()

	_, err = Db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			companyId INTEGER NOT NULL ,
			name TEXT NOT NULL,
			age INTEGER NOT NULL,
			isDeleted INTEGER DEFAULT 0
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

}

//Function to be used in different packages to initialise the connection to the db.

func ConnectToDb() {
	var err error
	Db, err = sql.Open("sqlite3", "./users.db")
	if err != nil {
		//Error connecting to db
		log.Fatal("Error connecting to db")
		log.Fatal(err)
	}
}

// Function to be used in different packages to close the connection to the db.
func CloseDb() {
	Db.Close()
}

//Helper to delete the particular row after testing

func DeleteRow(id int) {
	ConnectToDb()
	defer CloseDb()
	_, err := Db.Exec("DELETE FROM users WHERE companyId = ?", id)
	if err != nil {
		log.Fatal(err)
	}
}
