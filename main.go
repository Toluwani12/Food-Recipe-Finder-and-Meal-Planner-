package main

import (
	"Food/pkg/ingredient"
	"Food/pkg/recipe"
	"Food/pkg/users"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

func main() {
	r := chi.NewRouter()

	db, err := sqlx.Open("postgres", "postgres://olusolaalao:postgres@localhost:5432/recipe?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Check the connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	r.Mount("/recipe", recipe.NewResource(db).Router())

	r.Mount("/ingredient", ingredient.NewResource(db).Router())

	r.Mount("/users", users.NewResource(db).Router())

	log.Println("Server starting on port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
