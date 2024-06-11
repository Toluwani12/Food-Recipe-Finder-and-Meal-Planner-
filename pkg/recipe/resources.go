package recipe

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

	usrPrefRepo := user_preference.NewRepository(rs.db)
	usrPrefSvc := user_preference.NewService(usrPrefRepo)
	hndlr := NewHandler(svc, usrPrefSvc)

	r.Post("/search", hndlr.search)
	r.Group(func(r chi.Router) {
		r.Use(auth.AuthMiddleware)
		r.Get("/{id}/like", hndlr.like)
		r.Get("/{id}", hndlr.get)
		r.Get("/", hndlr.list)
		r.Post("/", hndlr.save)
		r.Delete("/{id}", hndlr.delete)
	})

	return r
}
