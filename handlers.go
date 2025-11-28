package main

import (
	_ "embed"
	"html/template"
	"net/http"
)

//go:embed html/home.html
var homeTmpl string

//go:embed html/login.html
var loginTmpl string

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	ts, err := template.New("home").Parse(homeTmpl)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "home", nil)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	ts, err := template.New("home").Parse(loginTmpl)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodGet {
		err = ts.ExecuteTemplate(w, "login", nil)
		if err != nil {

			app.errorLog.Println(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			app.errorLog.Println(err.Error())
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		email := r.PostForm.Get("email")
		password := r.PostForm.Get("password")

		// TODO: replace this with db lookup
		if email == "user@example.com" && password == "password123" {
			w.Header().Set("HX-Redirect", "/")
			w.WriteHeader(http.StatusOK)			
			return
		}
		
		w.Write([]byte(`<p style="color: red;">Invalid email or password</p>`))
	}
}
