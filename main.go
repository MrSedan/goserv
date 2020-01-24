package main

import (
	"net/http"

	common "./common"

	"github.com/gorilla/mux"
	"github.com/thedevsaddam/renderer"
)

var (
	router = mux.NewRouter()
	rnd    *renderer.Render
)

func init() {
	opts := renderer.Options{
		ParseGlobPattern: "./templates/*.html",
	}
	rnd = renderer.New(opts)
}

func main() {
	router.HandleFunc("/", common.LoginPageHandler)

	router.HandleFunc("/index", common.IndexPageHandler)
	router.HandleFunc("/login", common.LoginHandler).Methods("POST")

	router.HandleFunc("/register", common.RegisterPageHandler).Methods("GET")
	router.HandleFunc("/register", common.RegisterHandler).Methods("POST")

	router.HandleFunc("/logout", common.LogoutHandler).Methods("POST")

	http.Handle("/", router)
	http.ListenAndServe(":8000", nil)

}
