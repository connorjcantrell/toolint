package main

import (
	"log"
	"net/http"
	"os"

	"github.com/connorjcantrell/toolint/postgres"
	"github.com/connorjcantrell/toolint/web"
)

func main() {
	dsn := os.Getenv("TOOLINT_DB")
	store, err := postgres.NewStore(dsn)
	if err != nil {
		log.Fatal(err)
	}

	sessions, err := web.NewSessionManager(dsn)
	if err != nil {
		log.Fatal(err)
	}

	key := os.Getenv("CSRF_KEY")
	csrfKey := []byte(key)
	h := web.NewHandler(store, sessions, csrfKey)
	http.ListenAndServe(":3000", h)
}
