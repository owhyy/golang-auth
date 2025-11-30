package main

import (
	"log"
	"net/http"
	"os"
	"embed"
	
	"owhyy/simple-auth/models"
)

//go:embed static/*
var staticFS embed.FS

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	users    *models.UserModel
	tokens   *models.ValidationTokenModel
}

func main() {
	infoLog := log.New(os.Stdout, "INFO\t", log.LstdFlags)
	errorLog := log.New(os.Stderr, "ERROR\t", log.LstdFlags|log.Lshortfile)

	db, err := models.NewDB("./users.db")
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		users:    &models.UserModel{DB: db},
		tokens:    &models.ValidationTokenModel{DB: db},				
	}

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.FS(staticFS))
        mux.Handle("/static/", fileServer)	

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/login", app.login)
	mux.HandleFunc("/signup", app.signup)
	mux.HandleFunc("/verify", app.verify)

	srv := &http.Server{Addr: "0.0.0.0:8080", ErrorLog: errorLog, Handler: mux}
	infoLog.Println("Starting server on 0.0.0.0:8080")

	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
