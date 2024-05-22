package main

import (
	"Food/pkg/recipe"
	"Food/pkg/users"
	"context"
	"github.com/chromedp/chromedp"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"time"
)

func main() {
	r := chi.NewRouter()

	// Setting up CORS
	cors := cors.New(cors.Options{
		// Adjust settings based on your needs
		AllowedOrigins:   []string{"*"}, // or use `*` for allowing any origin
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not to need preflight request
	})

	// Use the CORS middleware
	r.Use(cors.Handler)

	db, err := sqlx.Open("postgres", "postgres://olusolaalao:postgres@localhost:5432/recipe?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// Check the connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	r.Mount("/recipes", recipe.NewResource(db).Router())

	//r.Mount("/ingredients", ingredient.NewResource(db).Router())

	r.Mount("/users", users.NewResource(db).Router())

	log.Println("Server starting on port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
