package user_preference

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

	r.Use(auth.AuthMiddleware)
	r.Post("/", hndlr.add)
	r.Delete("/", hndlr.delete)
	r.Put("/", hndlr.update)
	r.Get("/", hndlr.get)

	return r
}
