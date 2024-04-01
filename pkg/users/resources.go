package users

import (
	"Food/auth"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

type Resource struct {
	db *sqlx.DB
}

// NewResource creates and returns a resource.
func NewResource(db *sqlx.DB) *Resource {
	return &Resource{
		db: db,
	}
}

func (rs *Resource) Router() *chi.Mux {
	r := chi.NewRouter()

	repo := NewRepository(rs.db)
	svc := NewService(repo)
	hndlr := NewHandler(svc)

	r.Post("/login", hndlr.login)
	r.Post("/", hndlr.save)

	r.Group(func(r chi.Router) {
		r.Use(auth.AuthMiddleware)
		r.Get("/{id}", hndlr.get)
		r.Get("/", hndlr.list)
		r.Delete("/{id}", hndlr.delete)
		r.Put("/{id}", hndlr.update)
	})

	return r
}
