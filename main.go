package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is required")
	}

	var db *sql.DB
	for i := 0; i < 30; i++ {
		var err error
		db, err = sql.Open("postgres", dsn)
		if err != nil {
			log.Printf("sql.Open: %v", err)
			time.Sleep(1 * time.Second)
			continue
		}
		if err := db.Ping(); err != nil {
			log.Printf("ping: %v", err)
			_ = db.Close()
			db = nil
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}
	if db == nil {
		log.Fatal("could not connect to postgres")
	}
	defer db.Close()

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})

	http.HandleFunc("/dbz", func(w http.ResponseWriter, r *http.Request) {
		if err := db.Ping(); err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
		fmt.Fprintln(w, "db ok")
	})

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
