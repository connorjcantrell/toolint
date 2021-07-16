package main

import (
	"log"
	"net/http"

	"github.com/connorjcantrell/toolint/postgres"
	"github.com/connorjcantrell/toolint/web"
)

func main() {
	// dsn := os.Getenv("URL_DB")
	dsn := "postgres://root:secret@localhost/toolint?sslmode=disable"

	store, err := postgres.NewStore(dsn)
	if err != nil {
		log.Fatal(err)
	}

	sessions, err := web.NewSessionManager(dsn)
	if err != nil {
		log.Fatal(err)
	}

	csrfKey := []byte("01234567890123456789012345678901")
	h := web.NewHandler(store, sessions, csrfKey)
	http.ListenAndServe(":3000", h)
}
