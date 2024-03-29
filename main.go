package main

import (
	"Food/pkg/recipe"
	user2 "Food/pkg/user"
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

	user2.InitDB()

	r.Post("/register", user2.RegisterUser)

	log.Println("Server starting on port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
