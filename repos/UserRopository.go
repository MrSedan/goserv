package repos

import (
	"context"
	"database/sql"
	"log"

	"golang.org/x/crypto/bcrypt"

	//SQL-Driver
	_ "github.com/mattn/go-sqlite3"
)

var (
	ctx context.Context
	db  *sql.DB
)

//UserIsValid Checking User in DataBase on Login
func UserIsValid(uName, pwd string) bool {
	var err error
	db, err := sql.Open("sqlite3", "./serv.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	var passHash string
	row := db.QueryRow("SELECT passHash from users WHERE uName=$1", uName)
	err = row.Scan(&passHash)
	if err == sql.ErrNoRows {
		return false
	}
	if err != nil {
		log.Fatal(err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(passHash), []byte(pwd))
	if err != nil {
		return false
	}
	return true
}

//MaybeUser Checking for user in database
func MaybeUser(uName string) bool {
	var err error
	db, err := sql.Open("sqlite3", "./serv.db")
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

//Register is register users
func Register(uName, email, pwd string) {
	var err error
	db, err := sql.Open("sqlite3", "./serv.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	_, err = db.Exec("INSERT INTO users (uName, email, passHash) VALUES (?, ?, ?)", uName, email, HashAndSalt(pwd))
	if err != nil {
		log.Fatal(err)
	}
}

//HashAndSalt generates the Password Hash
func HashAndSalt(pwd string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		log.Fatal(err)
	}
	return string(hash)
}
