package repos

import (
	"database/sql"
	"log"

	"golang.org/x/crypto/bcrypt"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func UserIsValid(uName, pwd string) bool {
	var err error
	db, err = sql.Open("sqlite3", "../serv.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	quer, err := db.Query("SELECT passHash from users WHERE uName=\"?\"", uName)
	if err != nil {
		log.Fatal(err)
	}
	defer quer.Close()
	for quer.Next() {
		var passHash string
		err = quer.Scan(&passHash)
		if err != nil {
			log.Fatal(err)
		}
		err = bcrypt.CompareHashAndPassword([]byte(passHash), []byte(pwd))
		if err != nil {
			return false
		}
		return true
	}
	return false
}

func MaybeUser(uName string) bool {
	var err error
	db, err = sql.Open("sqlite3", "../serv.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM users WHERE uName=\"?\"", uName)
	defer rows.Close()
	for rows.Next() {
		return true
	}
	return false
}

func Register(uName, email, pwd string) {
	var err error
	db, err = sql.Open("sqlite3", "../serv,db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	_, err = db.Exec("INSERT INTO users (uName, email, passHash) VALUES (?, ?, ?)", uName, email, HashAndSalt(pwd))
	if err != nil {
		log.Fatal(err)
	}
}

func HashAndSalt(pwd string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		log.Fatal(err)
	}
	return string(hash)
}
