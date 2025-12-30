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
		log.Fatal(err)
	}

	log.Println("✅ Done")
}
