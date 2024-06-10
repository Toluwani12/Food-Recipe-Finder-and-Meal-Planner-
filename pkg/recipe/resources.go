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

	r.Use(auth.AuthMiddleware)

	//r.Put("/{id}", hndlr.update)
	r.Get("/{id}/like", hndlr.like)
	r.Get("/{id}", hndlr.get)
	r.Get("/", hndlr.list)
	r.Post("/search", hndlr.search)
	//r.Post("/search2", hndlr.search2)
	r.Group(func(r chi.Router) {
		r.Use(auth.AuthMiddleware)
		r.Post("/", hndlr.save)
		r.Delete("/{id}", hndlr.delete)
	})

	return r
}
