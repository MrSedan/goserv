package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	helpers "../helpers"
	repos "../repos"
	"github.com/gorilla/securecookie"
)

type ViewData struct {
	Name string
}

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	var body, _ = helpers.LoadFile("templates/login.html")
	fmt.Fprintf(w, body)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	name := r.FormValue("name")
	pass := r.FormValue("password")
	redirectTarget := "/"
	if !helpers.IsEmpty(name) && !helpers.IsEmpty(pass) {
		_userIsValid := repos.UserIsValid(name, pass)
		if _userIsValid {
			SetCookie(name, w)
			redirectTarget = "/index"
		} else {
			redirectTarget = "/register"
		}
	}
	http.Redirect(w, r, redirectTarget, 302)
}

func RegisterPageHandler(w http.ResponseWriter, r *http.Request) {
	var body, _ = helpers.LoadFile("templates/register.html")
	fmt.Fprintf(w, body)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	uName := r.FormValue("username")
	email := r.FormValue("email")
	pwd := r.FormValue("password")
	confirmPwd := r.FormValue("confirmPassword")
	redirectTarget := "/"
	_uName, _email, _pwd, _confirmPwd := false, false, false, false
	_uName = !helpers.IsEmpty(uName)
	_email = !helpers.IsEmpty(email)
	_pwd = !helpers.IsEmpty(pwd)
	_confirmPwd = !helpers.IsEmpty(confirmPwd)
	if _uName && _confirmPwd && _email && _pwd && _confirmPwd && pwd == confirmPwd {
		if !repos.MaybeUser(uName) {
			redirectTarget = "/index"
			repos.Register(uName, email, pwd)
		} else {
			redirectTarget = "/register"
		}
		http.Redirect(w, r, redirectTarget, 302)
	} else {
		http.Redirect(w, r, "/register", 302)
	}
}

func IndexPageHandler(w http.ResponseWriter, r *http.Request) {
	data := ViewData{}
	data.Name = GetUserName(r)
	if !helpers.IsEmpty(data.Name) {
		tmpl, err := template.ParseFiles("templates/index.html", "templates/base.html")
		if err != nil {
			log.Fatal(err)
		}
		tmpl.Execute(w, data)

	} else {
		http.Redirect(w, r, "/", 302)
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	ClearCookie(w)
	http.Redirect(w, r, "/", 302)
}

func SetCookie(userName string, w http.ResponseWriter) {
	value := map[string]string{
		"name": userName,
	}
	if encoded, err := cookieHandler.Encode("cookie", value); err == nil {
		cookie := &http.Cookie{
			Name:  "cookie",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(w, cookie)
	}
}

func ClearCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "cookie",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
}

func GetUserName(r *http.Request) (userName string) {
	if cookie, err := r.Cookie("cookie"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("cookie", cookie.Value, &cookieValue); err == nil {
			userName = cookieValue["name"]
		}
	}
	return userName
}
