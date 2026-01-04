package main

import (
	"flag"
	"log"
	"owhyy/simple-auth/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	var (
		dsn      = flag.String("dsn", "app.db", "SQLite database file")
		username = flag.String("username", "", "Username for admin user")
		email    = flag.String("email", "", "Email for admin user")
		password = flag.String("password", "", "Password for admin user")
	)
	flag.Parse()

	if *username == "" || *email == "" || *password == "" {
		log.Fatal("Username, email or password were not provided")
	}

	db, err := models.Migrate(*dsn)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	log.Println("🧑‍💻 Creating admin account...")
	if err := models.CreateAdmin(db, *username, *email, *password); err != nil {
		log.Fatal(err)
	}

	log.Println("✅ Done")
}
