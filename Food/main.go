package main

import (
	"Food/recipe"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

func main() {
	r := chi.NewRouter()

	db, err := sqlx.Open("postgres", "postgres://postgres:postgres@localhost:5432/recipe?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Check the connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	checkoutsResource := recipe.NewResource(db)
	r.Mount("/recipe", checkoutsResource.Router())

	// Add routes for CheckoutCreateHandler, Update, Delete

	http.ListenAndServe(":8080", r)
}
