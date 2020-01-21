package main

import (
	"database/sql"
	"net/http"
	"text/template"

	_ "github.com/mattn/go-sqlite3"
)

var database *sql.DB

type info struct {
	Name    string
	Message string
}
type ViewData struct {
	Information []info
	Title       string
}

func main() {
	db, err := sql.Open("sqlite3", "serv.db")
	if err != nil {
		panic(err)
	}
	database = db
	defer db.Close()

	http.HandleFunc("/", handler)
	http.HandleFunc("/create", testForm)

	http.ListenAndServe(":8000", nil)
}

//Standart Handler
func handler(w http.ResponseWriter, r *http.Request) {
	rows, err := database.Query("SELECT * FROM test")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var infos []info
	for rows.Next() {
		v := info{}
		err := rows.Scan(&v.Name, &v.Message)
		if err != nil {
			panic(err)
		}
		infos = append(infos, v)
	}
	tmpl, _ := template.ParseFiles("templates/index.html")
	data := ViewData{
		Title:       "Test",
		Information: infos,
	}
	tmpl.Execute(w, data)
}

func testForm(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			panic(err)
		}
		name := r.FormValue("name")
		message := r.FormValue("message")
		_, err = database.Exec("INSERT INTO test (name, message) VALUES (?, ?)", name, message)
		if err != nil {
			panic(err)
		}
		http.Redirect(w, r, "/", 301)
	} else {
		http.ServeFile(w, r, "templates/create.html")
	}
}
