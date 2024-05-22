package users

import (
	"Food/auth"
	"Food/pkg/user_preference"
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

		r.Get("/", hndlr.list)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", hndlr.get)
			r.Mount("/preferences", user_preference.NewResource(rs.db).Router())
			r.Delete("/", hndlr.delete)
			r.Put("/", hndlr.update)
		})

	})

	return r
}
