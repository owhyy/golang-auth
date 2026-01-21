package main

import (
	"flag"
	"log"
	"owhyy/simple-auth/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	var (
		dsn       = flag.String("dsn", "app.db", "SQLite database file")
		userCount = flag.Int("users", 10, "Number of users to generate")
		postCount = flag.Int("posts", 1000, "Number of posts to generate")
	)
	flag.Parse()

	db, err := models.Migrate(*dsn)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	log.Println("🌱 Populating database...")
	if err := models.Populate(db, *userCount, *postCount); err != nil {
		log.Fatal("Populating database failed with error: ", err)
	}
	log.Println("✅  Database populated successfully! Added ", *userCount, " users and ", *postCount, " posts")

	log.Println("🌱 Creating admin account")
	if err := models.CreateAdmin(db, "admin", "admin@example.com", "admin"); err != nil {
		log.Fatal("Creating admin account failed with error: ", err)
	}
	log.Println("✅ Admin account created successfully!")
	log.Println("✅ Done")
}
